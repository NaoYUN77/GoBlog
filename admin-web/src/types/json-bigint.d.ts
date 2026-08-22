declare module 'json-bigint' {
  interface JSONBigOptions {
    storeAsString?: boolean
    useNativeBigInt?: boolean
    alwaysParseAsBig?: boolean
  }
  interface JSONBig {
    parse(text: string): unknown
    stringify(value: unknown): string
  }
  function JSONBigFactory(options?: JSONBigOptions): JSONBig
  export = JSONBigFactory
}
