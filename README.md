# bigmac

Logs are better when you can identify the author. That is what HMAC is for. Bigmac provides a few io.Writer implementations that prepend HMAC and optionally secret identifiers (just in case you have several authors or want to rotate secrets).

## A simple case

In a simple case a system might only use one secret and write to a totally different stream if the secret is changed.

A SimpleSigner simply signs the []byte it recieves without identifying the secret or buffering for a complete line.

    sw := bigmac.NewSimpleSigner(os.Stdout, []byte("This is a demo secret"))
    simple := log.New(sw, ``, log.Flags())
    
    simple.Println("Everyone shares a secret.")
    simple.Println("Useful in very simple scenarios.")
    simple.Println("But long lived processes are going to need key rotation.")

The logged messages are prepended with a standard base64 encoded HMAC (sha256). No additional spacing is provided as the HMAC is a fixed length.

    vSAxlBRnSl0O+HsPLb0xdlbeelJSei6NzgInU4hufZA=2017/02/11 11:25:23 Everyone shares a secret.
    UnznaVS2reJ9WSuDYIbQnY7KEw9FRQcvskn+dI7rmkA=2017/02/11 11:25:23 Useful in very simple scenarios.
    nL7mXdqM7GRZ4h0I27FHlj3hZDJTgwUlxfwBQzxLtjo=2017/02/11 11:25:23 But long lived processes are going to need key rotation.

## Long Lived Processes and Identifiable Secrets

Longer lived services or other cases where secrets should be rotated regularly require both HMAC signatures and secret identification. Identifying the secret used to sign a message lets the system both combine messages from multiple authors and version the key individual authors use to sign messages.

An IdentifiedSigner prepends both the MAC and name of the secret used for signing to the []byte it receives. No buffering is performed.

    idw1 := bigmac.NewIdentifiedSigner(os.Stdout, "Author.1", []byte("This is a demo secret"))
    ident := log.New(idw1, ``, log.Flags())
    
    ident.Println("You can like totally trust that Author.1 created this entry.")
    ident.Println("You can be sure that nobody modified it or spoofed the identify.")
    
    idw2 := bigmac.NewIdentifiedSigner(os.Stdout, "Author.2", []byte("This is a demo secret"))
    ident.SetOutput(idw2)
    ident.Println("Even better, when you rotate and change the key version you can still read the whole log.")

The logged message is prepended with two single-space terminated tokens. The first is the HMAC encoded with standard base64. The second is the string name of the secret used. In the above code the first two log lines would use ````Author.1```` and the third ````Author.2````.

    bQ75DJUKHXiah3gXulPqKncEJO0F9UzAx75+LnrW+4w= Author.1 2017/02/11 11:25:23 You can like totally trust that Author.1 created this entry.
    NFyFBZSX+Q5jcfASmIdTVn0M71EWIIBWp1GciF9knHA= Author.1 2017/02/11 11:25:23 You can be sure that nobody modified it or spoofed the identify.
    pqA8TYOqfGU/YN0ztPRj1wG3qFuTvLTtGV/na0Q1wGQ= Author.2 2017/02/11 11:25:23 Even better, when you rotate and change the key version you can still read the whole log.
