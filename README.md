# codekub

Trilha local de prática de programação no terminal, no estilo de um juiz
online, com os problemas do [beecrowd](https://beecrowd.com). Você lê o
enunciado no navegador, resolve no seu editor, e o codekub julga com os
exemplos do enunciado e mostra o veredito ao vivo (AC, PE, WA, TLE, RE, CE).

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

Rode `codekub` para abrir o painel. Lá você navega pelas fases (Iniciante,
Ad-Hoc, Strings, Matemática…), abre um problema e usa as teclas:

| Tecla | Ação |
|---|---|
| `/` | busca um problema pelo número ou pelo título, de qualquer tela |
| `o` | abre o enunciado no beecrowd, no navegador |
| `e` | abre sua solução no editor |
| `s` | julga a solução |
| `w` | modo watch: julga a cada vez que você salvar |
| `?` | ajuda |

Os testes são os exemplos do enunciado, então passar neles não garante que a
solução seja aceita no beecrowd. Todas as fases ficam liberadas desde o início.

Os mesmos recursos existem como comandos, para quem prefere o terminal:

```sh
codekub next          # próximo problema
codekub show 1001     # exemplos e link do enunciado
codekub open 1001     # abre o enunciado no navegador
codekub test 1001     # julga sua solução
codekub watch 1001    # julga a cada vez que você salvar
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
```

Os problemas vieram do projeto
[potigol/beecrowd](https://github.com/potigol/beecrowd): título, categoria,
link e exemplos. Para criar os que faltarem (o script nunca sobrescreve um
problema que já existe na trilha):

```sh
git clone https://github.com/potigol/beecrowd ~/potigol-beecrowd
ruby tools/import_beecrowd.rb ~/potigol-beecrowd
```

Para dar a um problema enunciado e testes próprios, edite o `statement.md` e
crie um `reference.rb` e um `generator.rb` na pasta dele; depois rode
`ruby tools/build_tests.rb track/<fase>/<problema>` para gerar os testes.

A trilha (enunciados, metadados e testes) vai embutida no binário; as soluções
de referência e os geradores ficam só no repositório. Uma versão é publicada
criando uma tag `vX.Y.Z`: o GitHub Actions gera os binários e atualiza o
Homebrew.

## Estrutura

```
track/<fase>/<id>-<slug>/   statement.md  meta.json  tests/  (reference.rb  generator.rb)
tools/import_beecrowd.rb    cria os problemas a partir de um clone de potigol/beecrowd
tools/build_tests.rb        gera tests/ rodando reference.rb nas entradas do generator.rb
internal/judge              juiz: Runner (LocalRunner), comparação, vereditos
internal/track              leitura da trilha
internal/store              progresso em SQLite
internal/ui                 painel (Bubble Tea) e saída dos comandos
```

## Conteúdo e licenças

O código está sob a licença MIT (`LICENSE`). Os problemas são do beecrowd:
títulos, exemplos e enunciados pertencem ao beecrowd e aos seus autores, e o
codekub só aponta para eles. Os dados vêm do projeto potigol/beecrowd (MIT).
Veja `track/LICENSE.md`.
