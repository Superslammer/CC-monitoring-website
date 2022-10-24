local ws, side, power

function main()
    if not setup() then
        return
    end

    power = peripheral.wrap(side)
    
end

function setup()
    -- Getting the side the power source is on
    print("what side is the power on?")
    side = read()
    term.clear()
    term.setCursorPos(1,1)

    -- Getting the websocket
    local err
    local url = "wss://example.tweaked.cc/echo"
    
    ws, err = http.websocket(url)
    if ws then
        ws.send("ping")
        local msg = ws.receive()
        if msg == "ping" then
            
        end
    else
        print("error")
        print(err)
        return false
    end

    return true
end

main()