{
  inputs = {
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixos-25.11";
    };
    systems = {
      url = "github:nix-systems/default-linux";
    };
  };

  outputs = { 
    self,
    nixpkgs,
    systems,
  }: rec {
    checks = nixpkgs.lib.genAttrs (import systems) (system: let
      overlays = [ (import ./overlays/pjsip-overlay.nix) ];
      pkgs = import nixpkgs { inherit system; inherit overlays; };
    in {
      default = pkgs.stdenv.mkDerivation {
        name = "doorpix test";

        doCheck = true;
        src = ./.;

        # build-time deps
        nativeBuildInputs = with pkgs; [
          go
        ];

        # run-time deps
        buildInputs = with pkgs; [
          pkg-config
          pjsip
          openh264
          libvpx
          libyuv
          libv4l
          libopus
          bcg729
          SDL2
          libx11
          openssl
          alsa-lib
          gst_all_1.gstreamer
          gst_all_1.gst-plugins-base
          gst_all_1.gst-plugins-good
          gst_all_1.gst-plugins-bad
          gst_all_1.gst-plugins-ugly
          gst_all_1.gst-libav
          gst_all_1.gst-vaapi
        ];

        configurePhase = ''
          export GOCACHE=$TMPDIR/go-cache
          export GOMODCACHE=$TMPDIR/go-modcache
          export XDG_CACHE_HOME=$TMPDIR/xgd-cache
        '';       

        buildPhase = ''
          go build -o $out/doorpix cmd/doorpix/main.go
        '';

        checkPhase = ''
          go test -v ./...
        '';
      };
    });

    devShells = nixpkgs.lib.genAttrs (import systems) (system: let
      overlays = [ (import ./overlays/pjsip-overlay.nix) ];
      pkgs = import nixpkgs { inherit system; inherit overlays; };
    in {
      default = pkgs.mkShell {
        packages = with pkgs; [
          vim
          go-task
          go
          pkg-config
          pjsip
          openh264
          libvpx
          libyuv
          libv4l
          libopus
          bcg729
          SDL2
          libx11
          openssl
          alsa-lib
          gst_all_1.gstreamer
          gst_all_1.gst-plugins-base
          gst_all_1.gst-plugins-good
          gst_all_1.gst-plugins-bad
          gst_all_1.gst-plugins-ugly
          gst_all_1.gst-libav
          gst_all_1.gst-vaapi
        ];
      };
    });
  };
}
