package main

import fmt

func main() {
        in := bufio.NewReader(os.Stdin)
        for {
                line, err := in.ReadBytes('\n')

                if err == io.EOF {
                        return
                }

                if err != nil {
                        log.Fatalf("ReadBytes: %s", err)
                }

                parse(string(line))
        }
}

func parse (line string) {

//


}
