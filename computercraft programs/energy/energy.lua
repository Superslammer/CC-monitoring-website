local expect = require "cc.expect"
local expect = expect.expect
local apiURL, dataWaitTime
local getEnergyHandles, getEnergyConversion, getMaxEnergyHandles
local energyType, energyPeri
local computerID

function Main()
    Setup()

    local timerID = os.startTimer(5)
    while true do
        local event = { os.pullEvent() }

        if event[1] == "key" and event[2] == keys.f then
            print("Exiting program ...")
            return
        elseif event[1] == "timer" then
            local dateTime = GetCurrentTimeFormatted()
            local RF = GetEnergy()
            local response, request = SendData(dateTime, RF)
            if not response then
                print("Unable to insert data: " .. request)
            end

            timerID = os.startTimer(dataWaitTime)
        elseif event[1] == "peripheral_detach" then
            os.cancelTimer(timerID)
            print("Energy storage device diconnected, please reattach a energy storage device and press enter")
            io.read()
            Setup()
            timerID = os.startTimer(5)
        end
    end
end

function Setup()
    -- Get .env variables
    local envFile = io.open(".env", "r")
    local env = textutils.unserialise(envFile._handle.readAll())
    apiURL = env["apiURL"]
    dataWaitTime = env["dataWaitTime"]

    -- Close .env file
    envFile._handle.close()

    -- Get energy storage information
    local jsonFile = io.open("energyStorageInformation.json", "r")
    local decodedJSON = textutils.unserialiseJSON(jsonFile._handle.readAll())
    getEnergyHandles = decodedJSON["energyHandles"]
    getEnergyConversion = decodedJSON["energyConversion"]
    getMaxEnergyHandles = decodedJSON["maxEnergyHandles"]

    -- Close JSON file
    jsonFile._handle.close()

    -- Get energy peripheral
    -- Checking order: Bottom --> Top --> Back --> Front --> Right --> Left
    for key, _ in pairs(getEnergyHandles) do
        energyPeri = peripheral.find(key)
        if energyPeri ~= nil and key ~= nil then
            energyType = key
            break
        end
    end

    -- Check if supported energy storage device is connected to computer
    if energyType == nil then
        print("No supported energy storage device is connected to the computer")
        error()
    end

    -- Check if URL is valid
    local isValid, err = http.checkURL(apiURL)
    if not isValid then
        print("Malformed api URL: " .. err)
        error()
    end

    -- Check if URL is online
    local res, err = Ping()
    if not res then
        print("Connection error: " .. err)
        error()
    end

    -- Register computer as energy computer if not done already
    RegisterComputer()
end

function RegisterComputer()
    SetUpComputerID()
    local label = os.getComputerLabel()
    local computerInfo = { ["computerID"] = computerID, ["maxEnergy"] = GetMaxEnergy(), ["name"] = label }
    local body = textutils.serialiseJSON(computerInfo)
    local headers = { ["Content-type"] = "application/json; charset=UTF-8" }
    local resHandler, err, errRes = http.post(apiURL .. "energyComputer", body, headers)
    HandleHttpErr(err, errRes)

    local resJSON = resHandler.readAll()
    local res = textutils.unserialiseJSON(resJSON)

    local registeredError = "Computer already registered as a energy computer"
    if res["error"] and res["msg"] == registeredError then
        print(res["msg"])
        return
    elseif res["error"] then
        print(res["msg"])
        error()
    end

    print(res["msg"])
end

function SetUpComputerID()
    computerID = os.getComputerID()
end

function GetCurrentTimeFormatted()
    return os.date("!%F %T")
end

function GetEnergy()
    local energyVal = energyPeri[getEnergyHandles[energyType]]()
    local convVal = getEnergyConversion[energyType]
    return energyVal * convVal
end

function GetMaxEnergy()
    local maxVal = energyPeri[getMaxEnergyHandles[energyType]]()
    local convVal = getEnergyConversion[energyType]
    return maxVal * convVal
end

function SendData(dateTime, rf)
    expect(1, dateTime, "string")
    expect(2, rf, "number")

    local data = { ["data"] = { { ["dateTime"] = dateTime, ["RF"] = rf, ["computerID"] = computerID } } }
    local json = textutils.serialiseJSON(data)
    local headers = { ["Content-type"] = "application/json; charset=UTF-8" }

    local res, err, errRes = http.post(apiURL .. "energyData", json, headers)
    HandleHttpErr(err, errRes)

    if res.readAll() == "Data inserted" then
        return true
    else
        return false, json
    end
end

function HandleHttpErr(err, errRes)
    if err == nil then
        return
    end

    if err == "Could not connect" then
        print("Unable to reach api url: ", apiURL)
        error()
    end
    print(err)
    error("", 1)
end

function Ping()
    local res, err, errRes = http.get(apiURL)
    if res ~= nil then
        return true, nil
    end
    return false, err
end

Main()
