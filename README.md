# bigmac

Logs are better when you can identify the author. That is what HMAC is for. Bigmac provides a few io.Writer implementations that prepend HMAC and optionally secret identifiers (just in case you have several authors or want to rotate secrets).

## A simple case

A SimpleSigner simply signs the []byte it recieves without identifying the secret or buffering for a complete line.

    w := bigmac.NewSimpleSigner(os.Stdout, []byte("This is a demo secret"))
    log.SetOutput(w)
    
    log.Println("Hello, World!")

An IdentifiedSigner prepends both the MAC and name of the secret used for signing to the []byte it receives. No buffering is performed.

    w := bigmac.NewIdentifiedSigner(os.Stdout, "Author.1", []byte("This is Author's secret"))
    log.SetOutput(w)
    
    log.Println("Oh, look! An author identifiable message!")

    w := bigmac.NewIdentifiedSigner(os.Stdout, "Author.2", []byte("This is Author's new secret"))
    log.SetOutput(w)
    log.Println("This entry was signed with the updated key for Author.")


