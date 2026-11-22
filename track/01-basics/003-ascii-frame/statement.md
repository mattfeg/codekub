# Moldura

Antes das interfaces gráficas, os programas desenhavam caixas no terminal
usando só caracteres comuns. Vamos fazer uma moldura em volta do nome
`codekub!`.

## Entrada

Não há entrada.

## Saída

Imprima a moldura abaixo, com 5 linhas:

- a primeira e a última linha são a borda: um `+`, 12 hifens (`-`) e outro `+`;
- a segunda e a quarta linha são um `|`, 12 espaços e outro `|`;
- a linha do meio é um `|`, 2 espaços, o texto `codekub!`, 2 espaços e outro `|`.

## Exemplo

Saída:

```
+------------+
|            |
|  codekub!  |
|            |
+------------+
```

## Observação

Os espaços dentro da moldura fazem parte da resposta. Se faltar um, as bordas
da direita ficam desalinhadas e a saída não confere com a esperada.
