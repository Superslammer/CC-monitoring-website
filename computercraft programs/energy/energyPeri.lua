local M = {}
local peri, handle
local energyInfo = require "energyInfo"

function M.findPeripheral()
    if peri ~= nil then
        return peri
    end

    -- Load energy storage information
    energyInfo.loadEnergyInfo("energyStorageInformation.json")

    -- Checking order: Bottom --> Top --> Back --> Front --> Right --> Left
    for key, _ in pairs(energyInfo.getEnergyHandles()) do
        peri = peripheral.find(key)
        if peri ~= nil and key ~= nil then
            handle = key
            break
        end
    end

    -- Check if supported energy storage device is connected to computer
    if handle == nil then
        print("No supported energy storage device is connected to the computer")
        error()
    end

    return peri
end

function M.getEnergyStored()
    local energyVal = peri[energyInfo.getEnergyHandles()[handle]]()
    local convVal = energyInfo.getEnergyConversion()[handle]
    return energyVal * convVal
end

function M.getMaxEnergyStored()
    local maxVal = peri[energyInfo.getMaxEnergyHandles()[handle]]()
    local convVal = energyInfo.getEnergyConversion()[handle]
    return maxVal * convVal
end

return M