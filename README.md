# yokan
自作言語のインタプリタです。なんとか完成しました。

## 動かし方
`go run main.go` で動きます。

## 構文

### リテラル

```js
1
"abc\n"
```
整数リテラルと文字列リテラルがあります。

### 計算

```js
+123
-123
1+1
1-1
1*1
1/1
1==1
1!=1
1<1
1<=1
1>1
1>=1
```
これらの種類の計算ができます。

```js
1==1
"str"=="str"
true==true
null==null
```
`==`と`!=`はこれらの型に対応しています。

### 変数

```js
abc = 123
abc
```
使えます。

### 関数

```js
func = (args){
  assign = args
  assign
}
```
関数を定義するには無名関数を変数に代入します。
関数は複数の文を持てます。

```js
val = 123
func = (){ val }
```
関数の外側の値を取得できます。

```js
val = 123
func = (){ val = 456 }
val
```
結果は`123`となり、外側の値を変更することはできません。

### 組み込み

```js
true
false
null
```
これらはキーワードではなく、組み込みの変数となっています。これらの名前の変数を作ることもできますが、控えるべきでしょう。

```js
puts("Hello world!")
puts(123)
puts(true)
puts(null)
puts((){123})
```
puts関数はオブジェクトを出力できます。

```js
if(cond, true_expr, false_expr)
```
条件分岐はif関数を使います。

```js
if(cond, (){aaa}, (){bbb})()
```
true_exprとfalse_exprは必ず評価されるため、評価されたくない場合には関数を使います。

## 例

### Hello world!

```js
puts("Hello world!")
```

### FizzBuzz

```js
remainder=(n,d){n-(n/d)*d}
fizzbuzzii=(n){if(remainder(n,15)==0,(){puts("FizzBuzz\n")},(){if(remainder(n,3)==0,(){puts("Fizz\n")},(){if(remainder(n,5)==0,(){puts("Buzz\n")},(){puts(n)})()})()})()}
fizzbuzzi=(max,n){fizzbuzzii(n) if(max>n,(){fizzbuzzi(max,n+1)},(){})()}
fizzbuzz=(max){fizzbuzzi(max,1)}
fizzbuzz(15)
```
`go run main.go`では1行ずつで入力が切られてしまうので改行を削っています。

展開したものはこちらになります。
```js
remainder=(n,d){n-(n/d)*d}

fizzbuzzii=(n){
	if(
		remainder(n,15)==0,
		(){puts("fizzbuzz\n")},
		(){
			if(
				remainder(n,3)==0,
				(){puts("fizz\n")},
				(){
					if(
						remainder(n,5)==0,
						(){puts("buzz\n")},
						(){puts(n)}
					)()
				}
			)()
		}
	)()
}

fizzbuzzi=(max,n){
	fizzbuzzii(n)
	if(
		max>n,
		(){fizzbuzzi(max,n+1)},
		(){}
	)()
}
fizzbuzz=(max){
	fizzbuzzi(max,1)
}
fizzbuzz(15)
```
