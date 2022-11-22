local expect = require "cc.expect"
local expect = expect.expect
local apiURL
local getEnergyHandles, getEnergyConversion, getMaxEnergyHandles
local energyType, energyPeri

function main()
    setup()

    local timerID = os.startTimer(5)
    while true do
        local event = { os.pullEvent() }

        if event[1] == "key" and event[2] == keys.f then
            print("Exiting program ...")
            return
        elseif event[1] == "timer" then
            local dateTime = getCurrentTimeFormatted()
            local RF = getEnergy()
            local response, request = sendData(dateTime, RF)
            if not response then
                print("Unable to insert data: " .. request)
            end

            timerID = os.startTimer(5)
        elseif event[1] == "peripheral_detach" then
            os.cancelTimer(timerID)
            print("Energy storage device diconnected, please reattach a energy storage device and press enter")
            io.read()
            setup()
            timerID = os.startTimer()
        end
    end
end

function setup()
    -- Get .env variables
    local envFile = io.open(".env", "r")
    local env = textutils.unserialise(envFile._handle.read())
    apiURL = env["apiURL"]

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
    local res, err = ping()
    if not res then
        print("Connection error: " .. err)
        error()
    end

end

function setUpComputerID()
    local computerID = os.getComputerID()
end

function getCurrentTimeFormatted()
    return os.date("!%F %T")
end

function getEnergy()
    local energyVal = energyPeri[getEnergyHandles[energyType]]()
    local convVal = getEnergyConversion[energyType]
    return energyVal * convVal
end

function getMaxEnergy()
    local maxVal = energyPeri[getMaxEnergyHandles[energyType]]()
    local convVal = getEnergyConversion[energyType]
    return maxVal * convVal
end

function sendData(dateTime, rf)
    expect(1, dateTime, "string")
    expect(2, rf, "number")

    local data = { ["data"] = { { ["dateTime"] = dateTime, ["RF"] = rf } } }
    local json = textutils.serialiseJSON(data)
    local headers = { ["Content-type"] = "application/json; charset=UTF-8" }

    local res, err, errRes = http.post(apiURL, json, headers)
    handleHttpErr(err, errRes)

    if res.readAll() == "Data inserted" then
        return true
    else
        return false, json
    end
end

function handleHttpErr(err, errRes)
    if err == nil then
        return
    end

    if err == "Could not connect" then
        print("Unable to reach api url: ", apiURL)
        error()
    end
    print(err)
    error()
end

function ping()
    local res, err, errRes = http.get(apiURL)
    if res ~= nil then
        return true, nil
    end
    return false, err
end

main()
