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
}
