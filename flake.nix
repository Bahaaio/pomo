{
  description = "A simple Pomodoro timer TUI application written in Go.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      allSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems =
        f:
        nixpkgs.lib.genAttrs allSystems (
          system:
          f {
            pkgs = import nixpkgs { inherit system; };
          }
        );
    in
    {
      packages = forAllSystems (
        { pkgs }:
        {
          default = pkgs.buildGoModule {
            pname = "pomo";
            version = "1.2.1";

            src = ./.;

            vendorHash = "sha256-kbTYq4Xc86bcmNMhInq1rwYTbGRmu2TEXT2e7bqT5YY=";

            ldflags = [
              "-s"
              "-w"
            ];

            meta = with pkgs.lib; {
              description = "Customizable TUI Pomodoro timer with ASCII art, progress bar, desktop notifications, and productivity statistics";
              homepage = "https://github.com/Bahaaio/pomo";
              license = licenses.mit;
              platforms = platforms.linux ++ platforms.darwin;
              mainProgram = "pomo";
            };
          };
        }
      );
    };
}
