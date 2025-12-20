{
  inputs,
  lib,
  config,
  pkgs,
  ...
}: 
let
  username = "picric_acid";
  secrets = if builtins.pathExists ./secrets.nix
    then import ./secrets.nix
    else { gitName = "Default Name"; gitEmail = "default@example.com"; };
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

  home.packages = with pkgs; [
    claude-code
  ];

  programs.home-manager.enable = true;
  
  programs.neovim = {
    enable = true;
    
    plugins = with pkgs.vimPlugins; [
      nvim-treesitter.withAllGrammars
      base16-vim
    ];
  };

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
    settings = {
      user = {
	name = secrets.gitName;
	email = secrets.gitEmail;
      };
      core = {
        editor = "nvim";
      };
    };
  };
}
