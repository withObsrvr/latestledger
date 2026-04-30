{
  description = "latestledger.com — Stellar latest ledger dashboard";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forAllSystems = f:
        nixpkgs.lib.genAttrs systems (system:
          f (import nixpkgs { inherit system; })
        );
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gotools
            delve
            templ
            jq
            curl
            git
          ];

          shellHook = ''
            export CGO_ENABLED=0
            export PS1="(latestledger-flake) $PS1"

            echo ""
            echo "Latest Ledger dev shell"
            echo "  templ generate: templ generate"
            echo "  run server:     go run ./cmd/latestledger"
            echo "  test:           go test ./..."
            echo ""
          '';
        };
      });
    };
}
