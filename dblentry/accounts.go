package dblentry

import "fmt"

import "github.com/prataprc/goparsec"
import "github.com/prataprc/golog"

var inclusives = []string{
	"asset", "liability", "capital", "equity", "income", "expense",
}

type Account struct {
	name     string
	virtual  bool
	balanced bool
	balance  float64
	// from account directive
	note    string
	check   string
	assert  string
	eval    string
	defblns bool
}

func NewAccount(name string) *Account {
	acc := &Account{name: name}
	return acc
}

func (acc *Account) Yledger(db *Datastore) parsec.Parser {
	y := parsec.OrdChoice(
		func(nodes []parsec.ParsecNode) parsec.ParsecNode {
			t := nodes[0].(*parsec.Terminal)
			name := string(t.Value)
			switch t.Name {
			case "FULLACCNM":
				acc.name = name
				acc.virtual, acc.balanced = false, true
				return acc
			case "VFULLACCNM":
				acc.name = name[1 : len(name)-1]
				acc.virtual, acc.balanced = true, false
				return acc
			case "BFULLACCNM":
				acc.name = name[1 : len(name)-1]
				acc.virtual, acc.balanced = true, true
				return acc
			}
			panic(fmt.Errorf("unreachable code: terminal(%q)", t.Name))
		},
		ytok_accname, ytok_vaccname, ytok_baccname,
	)
	return y
}

func (acc *Account) SetOpeningbalance(amount float64) *Account {
	acc.balance = amount
	return acc
}

func (acc *Account) SetDirective(account *Account) *Account {
	acc.note = account.note
	acc.check = account.check
	acc.assert = account.assert
	acc.eval = account.eval
	return acc
}

func (acc *Account) Virtual() bool {
	return acc.virtual
}

func (acc *Account) Balanced() bool {
	return acc.balanced
}

func (acc *Account) Apply(db *Datastore, trans *Transaction, p *Posting) error {
	acc.balance += p.commodity.amount
	fmsg := "%v balance (from %v <%v>): %v\n"
	log.Debugf(fmsg, acc.name, trans.desc, p.commodity.amount, acc.balance)
	return nil
}

func (acc *Account) String() string {
	return fmt.Sprintf("%v", acc.name)
}