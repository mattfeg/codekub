# Caracteres especiais

Alguns caracteres têm significado especial dentro de um texto no código: as
aspas marcam onde o texto começa e termina, e a barra invertida (`\`) serve
para escrever caracteres que não dá para digitar direto. Para imprimi-los,
você vai precisar descobrir como sua linguagem os trata.

## Entrada

Não há entrada.

## Saída

Imprima exatamente as quatro linhas abaixo:

1. `Aspas duplas: "codekub"`
2. `Aspas simples: 'codekub'`
3. `Barra invertida: C:\codekub\solucoes`
4. `Porcentagem: 100%`

## Exemplo

Saída:

```
Aspas duplas: "codekub"
Aspas simples: 'codekub'
Barra invertida: C:\codekub\solucoes
Porcentagem: 100%
```

## Dica

Procure por *sequências de escape* na documentação da sua linguagem. Em muitas
delas, `\"` imprime aspas duplas e `\\` imprime uma única barra invertida.
