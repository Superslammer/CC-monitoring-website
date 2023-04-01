local expect = require "cc.expect"
expect = expect.expect
local env
local energyPeri = require "energyPeri"
local computer = require "computer"

function Main()
    Setup()

    -- Data collection loop
    local timerID = os.startTimer(env["dataWaitTime"])
    while true do
        local event = { os.pullEvent() }

        if event[1] == "key" and event[2] == keys.f then
            print("Exiting program...")
            return
        elseif event[1] == "timer" then
            computer.sendEnergyData()
        elseif event[1] == "key" and event[2] == keys.u then
            -- Update peripheral
            energyPeri.findPeripheral()
            computer.registerComputer()
        elseif evnet[1] == "peripheral_detach" then
            os.cancelTimer(timerID)
            print("Energy storage device diconnected, please reattach a energy storage device and press enter")
            io.read()
            Setup()
            timerID = os.startTimer(5)
        end
    end
end

function Setup()
    -- Read .env file
    local envFile = io.open(".env", "r")
    env = textutils.unserialise(envFile._handle.readAll())
    envFile._handle.close()

    -- Check if URL is valid
    local isValid, err = http.checkURL(env["apiURL"])
    if not isValid then
        print("Malformed api URL: " .. err)
        error()
    end

    -- Check if URL is online
    local res, err = computer.ping(env["apiURL"])
    if not res then
        print("Connection error: " .. err)
        error()
    end

    -- Register computer as energy computer if not done already
    computer.setup(env["apiURL"])
    computer.registerComputer()
end

Main()
