# codekub

Trilha local de prática de programação no terminal, no estilo de um juiz
online: do *Olá, mundo* aos problemas de competição. Você resolve no seu
editor, e o codekub julga com testes e mostra o veredito ao vivo (AC, PE, WA,
TLE, RE, CE).

*by Matheus Ferreira*

## Instalar

Você precisa da linguagem que vai praticar instalada: hoje, Ruby ou Python.

**macOS e Linux** (Homebrew):

```sh
brew install mattfeg/tap/codekub
```

**Windows:** baixe `codekub_<versão>_windows_amd64.zip` (ou `arm64`) em
[Releases](https://github.com/mattfeg/codekub/releases), extraia o
`codekub.exe` para uma pasta e coloque essa pasta no PATH. No PowerShell:

```powershell
$dir = "$env:LOCALAPPDATA\codekub"
New-Item -ItemType Directory -Force $dir | Out-Null
Expand-Archive codekub_*_windows_amd64.zip -DestinationPath $dir -Force
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$dir", "User")
```

Abra um terminal novo (de preferência o Windows Terminal) e rode
`codekub doctor`. Para Ruby no Windows, use o [RubyInstaller](https://rubyinstaller.org).

**Com Go instalado**, em qualquer sistema:

```sh
go install github.com/mattfeg/codekub@latest
```

Depois de instalar, `codekub doctor` confere linguagens, isolamento e pastas.

## Usar

```sh
codekub              # painel interativo
codekub next         # próximo exercício
codekub show 13      # enunciado
codekub test 13      # julga sua solução
codekub test 13 --exemplos
codekub watch 13     # julga a cada vez que você salvar
codekub progress
codekub help
```

Suas soluções ficam em `~/codekub/solucoes/` (criadas com um modelo quando você
abre um exercício) e o progresso em `~/.codekub/codekub.db`. Variáveis:
`CODEKUB_SOLUTIONS`, `CODEKUB_HOME`, `CODEKUB_LANG` e `CODEKUB_TRACK`.

## Configurar

```sh
codekub config                    # mostra as configurações e de onde vem cada uma
codekub config editor nvim        # ou nano, micro, code, "code -w", subl…
codekub config editor --padrao    # volta para $VISUAL, $EDITOR ou o do sistema
codekub config linguagem python
```

No painel, a aba **Ajustes** faz o mesmo: lista os editores instalados na sua
máquina (ou aceita um comando personalizado) e as linguagens disponíveis.

O editor é escolhido nesta ordem: `CODEKUB_EDITOR`, `~/.codekub/config.json`,
`$VISUAL`, `$EDITOR` e, por fim, o editor de texto padrão do sistema. Editores
de terminal abrem dentro do painel e voltam quando você sai deles. Editores
gráficos abrem em outra janela e combinam bem com o modo watch.

## Como as soluções rodam

Na sua máquina, com o interpretador instalado, numa pasta temporária, com
limite de tempo (por linguagem), limite de saída e encerramento de todos os
processos que a solução criar. No macOS, também dentro do `sandbox-exec`: sem
rede e sem escrita fora da pasta temporária. No Linux e no Windows não há
sandbox, e em nenhum sistema há limite de memória no modo local.

## Desenvolver

```sh
go build -o codekub .
CODEKUB_TRACK=./track ./codekub    # usa a trilha do disco em vez da embutida
ruby tools/build_tests.rb          # regenera os testes depois de editar exercícios
```

A trilha (enunciados, metadados e testes) vai embutida no binário; as soluções
de referência e os geradores ficam só no repositório. Uma versão é publicada
criando uma tag `vX.Y.Z`: o GitHub Actions gera os binários e atualiza o
Homebrew.

## Estrutura

```
track/<fase>/<id>-<slug>/   statement.md  meta.json  reference.rb  generator.rb  tests/
tools/build_tests.rb        gera tests/ rodando reference.rb nas entradas do generator.rb
internal/judge              juiz: Runner (LocalRunner), comparação, vereditos
internal/track              leitura da trilha
internal/store              progresso em SQLite
internal/ui                 painel (Bubble Tea) e saída dos comandos
```

## Conteúdo e licenças

O código está sob a licença MIT (`LICENSE`). Os exercícios originais são de
Matheus Ferreira, sob CC BY 4.0 (`track/LICENSE.md`). Os chefões são adaptados
do Codeforces via o dataset `open-r1/codeforces` (CC BY 4.0), com crédito no
enunciado e no `meta.json` de cada um.
