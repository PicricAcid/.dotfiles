{
  inputs,
  lib,
  config,
  pkgs,
  ...
}: let
  username = "picric_acid";
in {
  nixpkgs = {
    overlays = [
      inputs.neovim-nightly-overlay.overlays.default
    ];
    config = {
      allowUnfree = true;
    };
  };

  home = {
    username = username;
    homeDirectory = "/Users/${username}";
    
    stateVersion = "24.05";
  };

  programs.home-manager.enable = true;
  
  programs.neovim.enable = true;
  xdg.configFile."nvim" = {
    source = ../nvim;
    recursive = true;
  };

  programs.zsh = {
    enable = true;
    
    shellAliases = {
      vi = "nvim";
      la = "ls -a";
      ll = "ls -l";
    };
  };

  programs.starship = {
    enable = true;
    enableZshIntegration = true;
    settings = builtins.fromTOML (builtins.readFile ../starship.toml);
  };

  programs.wezterm.enable = true;
  xdg.configFile."wezterm" = {
    source = ../wezterm;
    recursive = true;
  };

  programs.git = {
    enable = true;
    userName = "PicricAcid";
    userEmail = "horumuarudehidohorumarin012@gmail.com";

    extraConfig = {
      core = {
        editor = "nvim";
      };
    };
  };
}
