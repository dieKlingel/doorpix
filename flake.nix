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
    devShells = nixpkgs.lib.genAttrs (import systems) (system: let
      pkgs = import nixpkgs { inherit system; };
    in {
      default = pkgs.mkShell {
        packages = with pkgs; [
          vim
          go-task
          go
          pkg-config
          pjsip
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
