package main

import "fmt"
import "reflect"

// GoWord a wrapper for functions that implement pieces of PISC
type GoWord func(*machine) error

// NilWord a wrapper for GoWords that should never fail
func NilWord(f func(*machine)) GoWord {
	return GoWord(func(m *machine) error {
		f(m)
		return nil
	})
}

func (m *machine) addGoWord(name, docstring string, impl GoWord) {
	m.helpDocs[name] = docstring
	m.predefinedWords[name] = impl
}

func t(m *machine) error {
	m.pushValue(Boolean(true))
	return nil
}

func f(m *machine) error {
	m.pushValue(Boolean(false))
	return nil
}

func dip(m *machine) error {
	quot := m.popValue().(*quotation).toCode()
	a := m.popValue()
	err := m.execute(quot)
	if err != nil {
		return err
	}
	m.pushValue(a)
	return nil
}

func pickDup(m *machine) error {
	distBack := int(m.popValue().(Integer))
	if distBack > len(m.values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.values))
	}
	m.pushValue(m.values[len(m.values)-distBack-1])
	return nil
}

func pickDrop(m *machine) error {
	distBack := int(m.popValue().(Integer))
	if distBack > len(m.values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.values))
	}
	valIdx := len(m.values) - distBack - 1
	val := m.values[valIdx]
	m.values = append(m.values[:valIdx], m.values[valIdx+1:]...)
	m.pushValue(val)
	return nil
}

func pickDel(m *machine) error {
	distBack := int(m.popValue().(Integer))
	if distBack > len(m.values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.values))
	}
	valIdx := len(m.values) - distBack - 1
	m.values = append(m.values[:valIdx], m.values[valIdx+1:]...)
	return nil
}

func lenEntry(m *machine) error {
	length := m.popValue().(lenable).Length()
	m.pushValue(Integer(length))
	return nil
}

func errorFromEntry(m *machine) error {
	msg := m.popValue().String()
	return fmt.Errorf(msg)
}

func reflectEq(m *machine) error {
	a := m.popValue()
	b := m.popValue()
	m.pushValue(Boolean(reflect.DeepEqual(a, b)))
	return nil
}

func (m *machine) loadPredefinedValues() {
	if m.predefinedWords == nil {
		panic("Uninitialized stack machine!")
	}
	m.addGoWord("t", "( -- t )", GoWord(t))
	m.addGoWord("f", "( -- f )", GoWord(f))
	m.addGoWord("dip", "( a quot -- ... a )", GoWord(dip))
	m.predefinedWords["pick-dup"] = GoWord(pickDup)
	m.predefinedWords["pick-drop"] = GoWord(pickDrop)
	m.predefinedWords["pick-del"] = GoWord(pickDel)
	m.addGoWord("len", "( e -- lenOfE ) ", GoWord(lenEntry))

	m.addGoWord("eq", " ( a b -- same? ) ", GoWord(runEq))
	// Discourage use of reflection based eq via long name
	m.addGoWord("deep-slow-reflect-eq", "( a b -- same? )", GoWord(reflectEq))
	m.addGoWord("error", "( msg -- !! )", GoWord(errorFromEntry))

	m.loadDebugWords()
	m.loadLocalWords()
	m.loadStringWords()
	m.loadBooleanWords()
	m.loadLoopWords()
	m.loadDictWords()
	m.loadVectorWords()
	m.loadSymbolWords()
	m.loadHigherMathWords()
	m.loadHelpWords()
	m.loadIOWords()
	m.loadShellWords()
	m.loadRandyWords()
	m.loadQuotWords()
	err := m.loadBoltWords()
	if err != nil {
		panic(fmt.Sprint("Error loading boltdb: ", err))
	}

	err = m.executeString(`"std_lib.pisc" import`, codePosition{source: "preload.go standard library import"})
	if err != nil {
		err = m.loadBackupPod()
		if err != nil {
			panic("Error loading pod! " + err.Error())
		}
	}
}

func (m *machine) loadBackupPod() error {
	podBackup := `:PRE -% ( dict? key -- dict ) [ ensure-dictionary ] dip [ t swap dict-set ] each-char ;
:PRE [] ( vec idx -- ) string>int vec-at ;
:PRE $: ( quot varName -- ) change ;
:PRE $ ( name -- val ) get-local ;
:PRE : ( val name -- ) set-local ;
:PRE -? ( dict key -- dict ? ) dict-has-key? ;
:PRE ++ ( varName -- . ) [ 1 + ] swap change ;
:PRE -- ( varname -- . ) [ 1 - ] swap change ;
:PRE -$ ( dict key -- val ) dict-get ;
:PRE -: ( dict? val key -- dict ) [ [ ensure-dictionary ] dip ] dip dict-set ;
:DOC 2vector ( a b -- vec )  ;
: 2vector ( a b -- vec ) <vector> 2 [ swap vec-prepend ] times ;
:DOC symb-eq ( symb symb -- eq? ) Are a and b equal symbols? ;
: symb-eq ( symb symb -- eq? ) symb-neq not ;
:DOC min ( a b -- smaller )  ;
: min ( a b -- smaller ) 2dup > [ nip ] [ drop ] if ;
:DOC over ( x y -- x y x )  ;
: over ( x y -- x y x ) 1 pick-dup ;
:DOC even? ( n -- ? )  ;
: even? ( n -- ? ) 2 mod zero? ;
:DOC quot>dict ( quot -- dict )  ;
: quot>dict ( quot -- dict ) get-locals call <dict> [ dict-set ] each-local drop-locals ;
:DOC show-locals ( -- ) Shows which locals are in current scope ;
: show-locals ( -- ) [ over >string " = " swap concat concat nip print ] each-local ;
:DOC nip ( a b -- b )  ;
: nip ( a b -- b ) 1 pick-del ;
:DOC abs ( a -- a )  ;
: abs ( a -- a ) dup 0 < [ -1 * ] when ;
:DOC 2dup ( a b -- a b a b ) Duplicates the top two elements of the stack ;
: 2dup ( a b -- a b a b ) 2 [ 1 pick-dup ] times ;
:DOC > ( a b -- ? )  ;
: > ( a b -- ? ) <= not ;
:DOC . ( a -- )  ;
: . ( a -- ) drop ;
:DOC = ( a b -- ? )  ;
: = ( a b -- ? ) - zero? ;
:DOC <= ( a b -- ? )  ;
: <= ( a b -- ? ) 2dup < [ = ] dip or ;
:DOC inspect ( a -- a )  ;
: inspect ( a -- a ) dup print ;
:DOC dict-if-not-dict ( .. -- dict )  ;
: dict-if-not-dict ( .. -- dict ) dup typeof "Dictionary" str-neq [ <dict> ] when ;
:DOC dup ( a -- a a ) Duplicates the top of the stack ;
: dup ( a -- a a ) 0 pick-dup ;
:DOC swap ( a b -- b a ) Swaps the top two elements of the stack ;
: swap ( a b -- b a ) 1 pick-drop ;
:DOC 2drop ( a b -- )  ;
: 2drop ( a b -- ) 2 [ drop ] times ;
:DOC keep ( ..a x quot: [ ..a x --- ..b ] -- ..b x )  ;
: keep ( ..a x quot: [ ..a x --- ..b ] -- ..b x ) over [ call ] dip ;
:DOC print ( a -- )  ;
: print ( a -- ) >string priv_puts ;
:DOC 3drop ( a b c -- )  ;
: 3drop ( a b c -- ) 3 [ drop ] times ;
:DOC when ( ? true -- res )  ;
: when ( ? true -- res ) [ ] ? call ;
:DOC >= ( a b -- ? )  ;
: >= ( a b -- ? ) < not ;
:DOC dict-if-empty-stack ( .. -- dict? )  ;
: dict-if-empty-stack ( .. -- dict? ) stack-empty? [ <dict> ] when ;
:DOC each2 ( v1 v2 quot -- v3 ) Applies a quotation to each elementwise pair in v1 and v2, to result in v3 ;
: each2 ( v1 v2 quot -- v3 ) get-locals :quot :v2 :v1 0 :i <vector> $v2 len $v1 len min [ $v1 $i vec-at $v2 $i vec-at quot vec-append [ 1 + ] $:i ] times drop-locals ;
:DOC bi ( a quot1 quot2 -- ... )  ;
: bi ( a quot1 quot2 -- ... ) [ keep ] dip call ;
:DOC str-neq ( str-a str-b -- eq? )  ;
: str-neq ( str-a str-b -- eq? ) str-eq not ;
:DOC dict-has-key? ( dict key -- dict bool )  ;
: dict-has-key? ( dict key -- dict bool ) "dict-has-key" extern-call ;
:DOC vec-reverse ( vec -- reversevec ) Reverses a vector ;
: vec-reverse ( vec -- reversevec ) get-locals dup :vec len :end 0 :i [ 1 - ] $:end $end 2 / [ $vec $i vec-at :x $vec $end vec-at :y $vec $y $i vec-set-at :vec $vec $x $end vec-set-at :vec [ 1 - ] $:end [ 1 + ] $:i ] times $vec drop-locals ;
:DOC drop ( a -- )  ;
: drop ( a -- ) 0 pick-del ;
:DOC quot>vector ( quot -- vec ) Take all the elements placed on the stack by quot, and put them in an array. Not local clean. ;
: quot>vector ( quot -- vec ) get-locals <vector> :vec <symbol> dup :mark /* Mark the stack */ swap call /* Fill the stack with info from the quotation */ [ dup $mark symb-neq ] [ [ swap vec-append ] $:vec ] while drop /* the mark */ $vec vec-reverse drop-locals ;
:DOC clear-stack ( -- )  ;
: clear-stack ( -- ) [ stack-empty? not ] [ drop ] while ;
:DOC divisor? ( n m -- ? )  ;
: divisor? ( n m -- ? ) mod zero? ;
:DOC ensure-dictionary ( .. -- dict )  ;
: ensure-dictionary ( .. -- dict ) dict-if-empty-stack dict-if-not-dict ;
:DOC splat ( vec -- items... ) Dump the contents of the vector onto the stack ;
: splat ( vec -- items... ) [ ] vec-each ;
:DOC max ( a b -- larger )  ;
: max ( a b -- larger ) 2dup > [ drop ] [ nip ] if ;
:DOC change ( quot varName -- .. )  ;
: change ( quot varName -- .. ) swap [ [ get-local ] keep ] dip dip set-local ;
:DOC if ( ? true false -- res ) if ? is t, call true quot, otherwise call falseQuot. Defined as '? call' ;
: if ( ? true false -- res ) ? call ;`
	return m.executeString(podBackup, codePosition{source: "Backup Pod"})
}
