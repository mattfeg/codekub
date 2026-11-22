# Tabela de preços

A padaria do bairro quer exibir o cardápio no terminal do caixa, em uma
tabela bem alinhada.

## Entrada

Não há entrada.

## Saída

Imprima a tabela abaixo. Ela tem duas colunas:

- **Produto**: o nome alinhado à esquerda, em um espaço de 13 caracteres;
- **Preço**: o valor com duas casas decimais, alinhado à direita, em um espaço
  de 7 caracteres.

Cada célula tem um espaço de cada lado entre as barras (`|`). As linhas de
borda usam `+` nos cruzamentos e `-` no resto. Os produtos, nesta ordem, são:

| Produto       | Preço |
|---------------|-------|
| Café          | 4.50  |
| Pão de queijo | 6.00  |
| Bolo de milho | 12.90 |
| Suco natural  | 8.75  |

## Exemplo

Saída:

```
+---------------+---------+
| Produto       |   Preço |
+---------------+---------+
| Café          |    4.50 |
| Pão de queijo |    6.00 |
| Bolo de milho |   12.90 |
| Suco natural  |    8.75 |
+---------------+---------+
```

## Observação

Repare que o cabeçalho `Preço` também fica alinhado à direita, como os
números.
