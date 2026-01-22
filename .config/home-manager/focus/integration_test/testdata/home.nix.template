{
  inputs,
  lib,
  config,
  pkgs,
  ...
}: {
  home = {
    username = "testuser";
    homeDirectory = "./test_workspace/testuser";
    stateVersion = "24.05";
  };

  home.packages = with pkgs; [
    git
  ];

  programs.home-manager.enable = true;
}
