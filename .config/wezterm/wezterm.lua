local wezterm = require 'wezterm'

local config = wezterm.config_builder()
config.automatically_reload_config = true

-- 日本語フォント設定
config.font = wezterm.font_with_fallback({
  'Menlo',
  'Monaco',
  'Hiragino Kaku Gothic ProN',
})

config.font_size = 14.0
config.use_ime = true
config.macos_forward_to_ime_modifier_mask = 'SHIFT|CTRL'

-- 日本語文字幅の調整
config.unicode_version = 14
config.allow_square_glyphs_to_overflow_width = "WhenFollowedBySpace"
config.window_background_opacity = 0.85
config.macos_window_background_blur = 20
config.hide_tab_bar_if_only_one_tab = true
config.color_scheme = 'Gogh (Gogh)'
--local mux = wezterm.mux
--wezterm.on("gui-startup", function(cmd)
--    local tab, pane, window = mux.spawn_window(cmd or {})
--    window:gui_window():toggle_fullscreen()
--end)
--
config.default_cursor_style = 'BlinkingBar'

-- Leader key設定
config.leader = { key = 'w', mods = 'CTRL', timeout_milliseconds = 1000 }

config.keys = {
  {key="Enter", mods="SHIFT", action=wezterm.action{SendString="\x1b\r"}},
  -- 画面分割
  {key="h", mods="LEADER", action=wezterm.action.SplitVertical{domain="CurrentPaneDomain"}},
  {key="v", mods="LEADER", action=wezterm.action.SplitHorizontal{domain="CurrentPaneDomain"}},
  -- ペインを閉じる
  {key="q", mods="LEADER", action=wezterm.action.CloseCurrentPane{confirm=true}},
  -- ペイン切り替え
  {key="LeftArrow", mods="LEADER", action=wezterm.action.ActivatePaneDirection("Left")},
  {key="RightArrow", mods="LEADER", action=wezterm.action.ActivatePaneDirection("Right")},
  {key="UpArrow", mods="LEADER", action=wezterm.action.ActivatePaneDirection("Up")},
  {key="DownArrow", mods="LEADER", action=wezterm.action.ActivatePaneDirection("Down")},
}


return config

