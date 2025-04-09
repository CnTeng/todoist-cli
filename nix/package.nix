{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "todoist-cli";
  version = "unstable-2025-04-09";

  src = ../.;

  vendorHash = "sha256-mTQL1i+V8YV8ePnowg4wLPhgCuyRBBLOK4bOm+s3pwo=";

  postFixup = ''
    mv $out/bin/todoist-cli $out/bin/todoist
  '';

  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    homepage = "https://github.com/CnTeng/todoist-cli";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [ CnTeng ];
    mainProgram = "todoist";
  };
}
