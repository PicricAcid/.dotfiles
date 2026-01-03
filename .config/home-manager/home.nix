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
  imports = [
    ./skill-hunter-packages.nix
  ];

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
    (buildGoModule {
      pname = "skill_hunter";
      version = "1.0.0";
      src = ./skill_hunter;   
      vendorHash = "lib.fakeHash"; 
    })
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

  xdg.configFile."wezterm" = {
    source = ../wezterm;
    recursive = true;
  };

  xdg.configFile."ghostty" = {
    source = ../ghostty;
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
      url = {
        "git@github.com:" = {
	  instedOf = "github.com";
	};
      };
    };
  };

  programs.ssh = {
    enable = true;
    enableDefaultConfig = false;

    matchBlocks = {
      "github.com" = {
        hostname = "github.com";
	user = "git";
	identityFile = "~/.ssh/id_ed25519";
      };
    };
  };

  programs.go.enable = true;
  home.sessionPath = [ "$HOME/go/bin" ];
}
