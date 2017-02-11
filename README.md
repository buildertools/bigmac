# bigmac

Logs are better when you can identify the author. That is what HMAC is for. Bigmac provides a few io.Writer implementations that prepend HMAC and optionally secret identifiers (just in case you have several authors or want to rotate secrets).

## A simple case

In a simple case a system might only use one secret and write to a totally different stream if the secret is changed.

A SimpleSigner simply signs the []byte it recieves without identifying the secret or buffering for a complete line.

    w := bigmac.NewSimpleSigner(os.Stdout, []byte("This is a demo secret"))
    log.SetOutput(w)
    
    log.Println("Hello, World!")

The logged message is prepended with a standard base64 encoded HMAC (sha256). No additional spacing is provided as the HMAC is a fixed length.

## Long Lived Processes and Identifiable Secrets

Longer lived services or other cases where secrets should be rotated regularly require both HMAC signatures and secret identification. Identifying the secret used to sign a message lets the system both combine messages from multiple authors and version the key individual authors use to sign messages.

An IdentifiedSigner prepends both the MAC and name of the secret used for signing to the []byte it receives. No buffering is performed.

    w := bigmac.NewIdentifiedSigner(os.Stdout, "Author.1", []byte("This is Author's secret"))
    log.SetOutput(w)
    
    log.Println("Oh, look! An author identifiable message!")

    w := bigmac.NewIdentifiedSigner(os.Stdout, "Author.2", []byte("This is Author's new secret"))
    log.SetOutput(w)
    log.Println("This entry was signed with the updated key for Author.")

The logged message is prepended with two single-space terminated tokens. The first is the HMAC encoded with standard base64. The second is the string name of the secret used. In the above code the first log line would use ````Author.1```` and the second ````Author.2````.
