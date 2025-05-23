{
  lib,
  buildGoModule,
  installShellFiles,
}:
buildGoModule {
  pname = "todoist-cli";
  version = "unstable-2025-05-23";

  src = ../.;

  vendorHash = "sha256-igt/FwxIZ+wvAM1vLuMcYFT0xiWZBB/eR97+aoF4xIA=";

  nativeBuildInputs = [ installShellFiles ];

  ldflags = [
    "-s"
    "-w"
  ];

  postInstall = ''
    mv $out/bin/todoist-cli $out/bin/todoist

    installShellCompletion --cmd todoist \
      --bash <($out/bin/todoist completion bash) \
      --zsh <($out/bin/todoist completion zsh) \
      --fish <($out/bin/todoist completion fish)
  '';

  meta = {
    homepage = "https://github.com/CnTeng/todoist-cli";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [ CnTeng ];
    mainProgram = "todoist";
  };
}
