local M = {}
local energyInfo

function M.loadEnergyInfo(filePath)
    local jsonFile = io.open(filePath, "r")._handle
    energyInfo = textutils.unserialiseJSON(jsonFile.readAll())
    jsonFile.close()
end

function M.getEnergyHandles()
    assert(type(energyInfo) ~= "nil", "loadEnergyInfo not called")
    return energyInfo["energyHandles"]
end

function M.getEnergyConversion()
    assert(type(energyInfo) ~= "nil", "loadEnergyInfo not called")
    return energyInfo["energyConversion"]
end

function M.getMaxEnergyHandles()
    assert(type(energyInfo) ~= "nil", "loadEnergyInfo not called")
    return energyInfo["maxEnergyHandles"]
end

return M