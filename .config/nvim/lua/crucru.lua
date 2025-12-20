local emoji = { "( -_-)", "( =_=).", "( -_-)..", "( =_=)..." } 
local spinner = { "/", "-", "\\", "|" }
local spinner_index = 1
local timer = nil

function start_spinner(mode)
    if timer then return end -- 既に回ってるなら何もしない
    
    if mode == "Emoji" then
	timer = vim.fn.timer_start(500, function()
	    vim.schedule(function()
    		vim.cmd("echohl ModeMsg")
		vim.cmd(string.format("echon 'Processing %s'", emoji[spinner_index]))
		vim.cmd("echohl None")
		spinner_index = spinner_index % #emoji + 1
	    end)
	end, { ["repeat"] = -1 })
    else
	timer = vim.fn.timer_start(100, function()
	    vim.schedule(function()
		vim.cmd("echohl ModeMsg")
		vim.cmd(string.format("echon '%s Processing...'", spinner[spinner_index]))
		vim.cmd("echohl None")
		spinner_index = spinner_index % #spinner + 1
	    end)
	end, { ["repeat"] = -1 })
    end
end

function stop_spinner()
    if timer then
	vim.fn.timer_stop(timer)
	timer = nil
	vim.cmd("echo ''")
    end
end

function on_exit(job_id, exit_code, event_type)
    vim.schedule(function()
	stop_spinner()
	print("done!")
    end)
end

function run_async_task()
    start_spinner("Emoji")
    vim.fn.jobstart({ "sleep", "10" }, {
	on_exit = on_exit
    })
end

vim.api.nvim_create_user_command("CrucruTest", run_async_task, {})

