# bigmac

Logs are better when you can identify the author. That is what HMAC is for. Bigmac provides a few io.Writer implementations that prepend HMAC and optionally secret identifiers (just in case you have several authors or want to rotate secrets). This package can be used just about anywhere since the builtin log package and most common logging packages (such as [Logrus](https://github.com/sirupsen/logrus)) log to io.Writer types.

## A Simple Case

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

## Asymmetric Identities

Sometimes (usually) your system needs to be able to verify the authenticity of a message without sharing a secret with the author. In those cases it is best to make use of ECDSA or RSA signatures. The ````IdentifiedPKCS1v15Signer```` and ````IdentifiedECDSASigner```` types provide such tooling.

### ECDSA

    ecpk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    ecdsa1 := bigmac.NewIdentifiedECDSASigner(os.Stdout, "generated-ECDSA", ecpk)
    ident.SetOutput(ecdsa1)
    ident.Println("This uses an ECDSA signature. Pretty fancy.")

Since ECDSA signatures produce two artifacts the log line is prepended with the R and S components (in that order) using standard base64 encoding. The secret name used for signing follows. All components are space delimited.

    Zd1nA9bknXGNdHoXT2e7Pkx/nC7M2rc3QXCzTk6qdZQ= x3LOGbVZmjnTvJG9T6qga6x7rRDsV7XzwiiNKmrlLBQ= generated-ECDSA 2017/02/11 15:44:27 This uses an ECDSA signature. Pretty fancy.

### PKCS1

    pk, _ := rsa.GenerateKey(rand.Reader, 2048)
    pkcs1 := bigmac.NewIdentifiedPKCS1v15Signer(os.Stdout, "generated-rsa", pk)
    ident.SetOutput(pkcs1)
    ident.Println("This uses PKCS1v15 with a 2048 bit key. The resulting signature is huge and takes \"forever\" to generate.")

The PKCS1 signature is standard base64 encoded and prepended to the provided input. The identity of the secret follows the signature. All components are space delimited.


