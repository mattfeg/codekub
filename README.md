# codekub

Trilha local de prática de programação no terminal, no estilo de um juiz
online: do *Olá, mundo* aos problemas de competição. Você resolve no seu
editor, e o codekub julga com testes e mostra o veredito ao vivo (AC, PE, WA,
TLE, RE, CE).

*by Matheus Ferreira*

## Instalar

Requer Go 1.22+ e a linguagem que você vai praticar (hoje: Ruby ou Python).

```sh
go build -o codekub .
./codekub doctor   # confere linguagens, isolamento e pastas
```

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

## Como as soluções rodam

Na sua máquina, com o interpretador instalado, numa pasta temporária, com
limite de tempo (por linguagem), limite de saída e um ambiente mínimo. No
macOS, também dentro do `sandbox-exec`: sem rede e sem escrita fora da pasta
temporária. Não há limite de memória no modo local.

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

Os exercícios originais são de Matheus Ferreira, sob CC BY 4.0. Os chefões
são adaptados do Codeforces via o dataset `open-r1/codeforces` (CC BY 4.0), com
crédito no enunciado e no `meta.json` de cada um.
