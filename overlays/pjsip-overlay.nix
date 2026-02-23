# pjsip-overlay.nix
self: super: {
  # Override the pjsip package with your custom configure flags
  pjsip = super.pjsip.overrideAttrs (oldAttrs: {
    configureFlags = (oldAttrs.configureFlags or []) ++ [
    ];

    buildInputs = oldAttrs.buildInputs ++ [
      super.pkg-config
      super.openh264
      super.libvpx
      super.libyuv
      super.libv4l
      super.libopus
      super.bcg729
      super.SDL2
      super.libx11
    ];

    # Optionally, add patches or extra files
    postPatch = ''
      ${oldAttrs.postPatch or ""}
      # Your custom patching commands here
      cp ${./config_site.h} pjlib/include/pj/config_site.h
    '';
  });
}

