{
  description = "Flair - A theme generation tool";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.golang-shared-configs.url = "github:curtbushko/golang-shared-configs";

  outputs = { self, nixpkgs, golang-shared-configs }:
    let
      goVersion = 25; # Change this to update the whole stack

      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        inherit system;
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
          overlays = [ self.overlays.default ];
        };
      });

      # Build go-ai-lint from source
      go-ai-lint = { pkgs }: pkgs.buildGoModule {
        pname = "go-ai-lint";
        version = "1.0.0";
        src = pkgs.fetchFromGitHub {
          owner = "curtbushko";
          repo = "go-ai-lint";
          rev = "v1.0.0";
          sha256 = "sha256-y2G7dTZqM/rEQaALu54bHigBeO1xxRIblBJ7QxOffW4=";
        };
        subPackages = [ "cmd/go-ai-lint" ];
        vendorHash = "sha256-zkXyXTEnMmBZnvzoq0UWKgzWZlyNRyQZCYAv+huZo0I=";
      };

      # Build flair from source
      flair = { pkgs }: pkgs.buildGoModule {
        pname = "flair";
        version = "0.1.0";
        src = ./.;
        subPackages = [ "cmd/flair" ];
        vendorHash = "sha256-5MZludg0yiSXVEj56AJ+rfak9bLQYJeCbPunlS9mX6A=";

        meta = with pkgs.lib; {
          description = "A theme generation tool for creating consistent color schemes";
          homepage = "https://github.com/curtbushko/flair";
          license = licenses.mit;
          maintainers = [];
        };
      };
    in
    {
      overlays.default = final: prev: {
        go = final."go_1_${toString goVersion}";
        flair = flair { pkgs = final; };
      };

      # Packages output for use as a flake input
      packages = forEachSupportedSystem ({ pkgs, system }: {
        flair = flair { inherit pkgs; };
        default = flair { inherit pkgs; };
      });

      devShells = forEachSupportedSystem ({ pkgs, system }:
        let
          sharedConfigs = golang-shared-configs.packages.${system}.all-configs;
        in {
        default = pkgs.mkShell {
          packages = with pkgs; [
            docker
            # go (version is specified by overlay)
            go
            go-task
            gotools
            golangci-lint
            (go-ai-lint { inherit pkgs; })
            sharedConfigs
          ];

          shellHook = ''
            cp -f ${sharedConfigs}/.golangci.yml .golangci.yml
            cp -f ${sharedConfigs}/.go-arch-lint.yml .go-arch-lint.yml
            cp -f ${sharedConfigs}/.go-ai-lint.yml .go-ai-lint.yml
          '';
        };
      });
    };
}
