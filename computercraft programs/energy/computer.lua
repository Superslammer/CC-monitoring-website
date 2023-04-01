local M = {}
local computerID
local apiURL
local energyPeri = require "energyPeri"
local name

function M.setup(URL)
    -- Get energy peripheral
    energyPeri.findPeripheral()

    computerID = os.getComputerID()
    apiURL = URL
    name = os.getComputerLabel()
end

local function getComputer(id)
    local res, err, errRes = http.get(apiURL .. "energy-computer/" .. id .. "/")
    if res == nil then
        print(err)
        return nil, false
    end

    local data = res.readAll()
    return textutils.unserialiseJSON(data), true
end

local function handleHttpCodes(code)
    assert(type(code) == "number", "HTTP code must be a number")
    if code == 404 then
        error("Resource unavaiable", 1)
    elseif code == 400 then
        error("Malformed request", 1)
    elseif code == 500 then
        error("Internal server error", 1)
    else
        error("HTTP code: " .. code, 1)
    end
end

function M.registerComputer()
    local data, isRegistered = getComputer(computerID)
    if isRegistered == true then
        name = data.name
        os.setComputerLabel(name)
        local body = {
            ["name"] = name,
            ["maxRF"] = energyPeri.getMaxEnergyStored()
        }
        M.updateComputer(body)
    else
        local body = {
            ["id"] = computerID,
            ["name"] = name,
            ["maxRF"] = energyPeri.getMaxEnergyStored(),
            ["currentRF"] = energyPeri.getEnergyStored()
        }
        local headers = { ["Content-type"] = "application/json; charset=UTF-8" }
        local method = "POST"
        local request = {
            url = apiURL .. "energy-computer/",
            body = textutils.unserialiseJSON(body),
            headers = headers,
            method = method
        }
        local _, err, errRes = http.post(request)
        if err ~= nil then
            local resCode = errRes.getResponseCode()
            handleHttpCodes(resCode)
        end
    end
end

function M.updateComputer(body)
    local headers = { ["Content-type"] = "application/json; charset=UTF-8" }
    local method = "UPDATE"
    local request = {
        url = apiURL .. "energy-computer/",
        body = textutils.unserialiseJSON(body),
        headers = headers,
        method = method
    }
    local _, err, errRes = http.post(request)
    if err ~= nil then
        local resCode = errRes.getResponseCode()
        handleHttpCodes(resCode)
    end
end

local function GetCurrentTimeFormatted()
    return os.date("!%F %T")
end

function M.sendEnergyData()
    local headers = { ["Content-type"] = "application/json; charset=UTF-8" }
    local method = "POST"
    local dateTime = GetCurrentTimeFormatted()
    local body = {
        ["dateTime"] = dateTime,
        ["RF"] = energyPeri.getEnergyStored(),
        ["computerID"] = computerID
    }
    local request = {
        url = apiURL .. "energy-data/",
        body = body,
        headers = headers,
        method = method
    }
    local _, err, errRes = http.post(request)
    if err ~= nil then
        local resCode = errRes.getResponseCode()
        handleHttpCodes(resCode)
    end
end

function M.ping(apiURL)
    local res, err = http.get(apiURL)
    if res ~= nil then
        return true, nil
    end
    return false, err
end

return M
