package ITP

import (
	"bufio"
	"net"
	"strconv"
)

type RGB struct {
	R byte
	G byte
	B byte
}

type image struct {
	height int
	width  int
	img    map[int]map[int]RGB
}

type itpConn struct {
	conn net.Conn
	rbuf *bufio.Reader
	wbuf *bufio.Writer
	req  string
}

type itpImage struct {
	conn itpConn
	img  image
}

func (this *itpImage) init() {
	this.conn.rbuf = bufio.NewReader(this.conn.conn)
	this.conn.wbuf = bufio.NewWriter(this.conn.conn)
}

func (this *itpImage) send() {
	for h := 0; h < this.img.height; h++ {
		for w := 0; w < this.img.width; w++ {
			color := this.img.img[h][w]
			line := string(color.R) + ";" + string(color.G) + ";" + string(color.B) + "\n"
			this.conn.req = this.conn.req + line
		}
	}
	this.conn.wbuf.WriteString(this.conn.req)
}

func (this *itpImage) read() {
	rd := this.conn.rbuf
	for h := 0; h < this.img.height; h++ {
		for w := 0; w < this.img.width; w++ {
			R, rerr := rd.ReadString(';')
			if rerr != nil {
				return // clear this.image
			}
			G, gerr := rd.ReadString(';')
			if gerr != nil {
				return // clear this.image
			}
			B, berr := rd.ReadString('\n')
			if berr != nil {
				return // clear this.image
			}
			r, _ := strconv.Atoi(R)
			g, _ := strconv.Atoi(G)
			b, _ := strconv.Atoi(B)
			this.img.img[h][w] = RGB{byte(r), byte(g), byte(b)}
		}
	}

}

type itpVideo struct {
	conn   itpConn
	frames map[int64]itpImage
}
