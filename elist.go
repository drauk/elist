// src/go/elist.go   2017-9-13   Alan U. Kennington.
// $Id: elist.go 46584 2017-09-09 02:22:08Z akenning $
// Error-message-stack class for error-message traceback.
// Using Go version go1.1.2.
/*-------------------------------------------------------------------------
Functions in this package.

Elist::
Elist::Error
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
New
Newf
Push
Pushf
-------------------------------------------------------------------------*/

/*
The "elist" Go-package implements an error-message-stack class "Elist", which
implements the standard "error" interface.

An Elist permits chaining of multiple error messages by pushing them onto a
single error-message-stack.
    type Elist struct {
        next    *Elist;         // Next node in a singly linked stack.
        value   interface{};    // The payload of the error node.
    }
The value-field is either a "string" or an "error", and this "error" may be any
struct which implements the standard "error" interface. The purpose of the
"error" option here is to permit non-Elist errors to be chained into an Elist
stack.

The next-field optionally points to the "Elist" structure which was returned by
a function-call which prompted the construction of this "Elist" structure. Thus
the next-field refers to an error indication which is chronologically prior this
error indication.

Usage example:
    func function0() error {
        var E error = function1();
        if E != nil {
            return elist.Push(E, "function0: error returned by function1");
        }
        return nil;
    }
Any kind of error interface-implementation can be chained into an Elist. The
output from Elist::Error() then has the form indicated in this example.
    Error 1: "function0: error returned by function1".
    Error 2: "function1: error returned by function2".
    Error 3: "function2: [error description etc.]".
When the calling function function0() does not receive an error message from a
called function function1(), a new Elist error-message-stack may be created as
in the following example.
    func function0() error {
        var n int = function1();
        if n < 0 {
            return elist.Newf(
                "function0: negative value returned by function1: %d", n);
        }
        return nil;
    }
Then the caller of function0() may either print the error message which is
returned, or else push a new message onto the returned message and pass that as
a return value.
*/
package elist

// External libraries.
import "fmt"

// import "strings"
// import "net/http"
// import "log"
// import "io"
// import "time"
// import "errors"

//=============================================================================
//=============================================================================

/*
Elist adds error-message-stack functions to the standard "error" interface. In
other words, "Elist" is a derived class of the standard class "error".
    type Elist struct {
        next    *Elist;         // Next node in a singly linked stack.
        value   interface{};    // The payload of the error node.
    }
The next-field permits the chaining of error messages so that calling functions
can print a trace-back of a sequence of errors instead of just the error message
of the last function which exits with an error status.
*/
type Elist struct {
    //------------------//
    //      Elist::     //
    //------------------//
    next  *Elist      // Next node in a singly linked stack.
    value interface{} // The payload of the error node.
}

/*
Return a plain-text error-message-stack, one error message per newline-separated
line.
Example output:
    Error 1: "handler_info(): args.fetch() failed".
    Error 2: "mj_args::fetch: 2 settings for parameter "i"".
Error messages are listed in LIFO order.
In other words, it's a trace-back from the outermost context to the innermost
context.
*/
func (p *Elist) Error() string {
    //------------------//
    //   Elist::Error   //
    //------------------//
    if p == nil {
        return ""
    }
    var str string = ""
    var n int = 1
    for q := p; q != nil; q = q.next {
        var msg string = ""

        switch x := q.value.(type) {
        case string:
            // An Elist object.
            msg = fmt.Sprintf(": \"%s\"", x)
        case error:
            // An old-style error object.
            //            msg = fmt.Sprintf(": \"%v\"", x.Error());
            msg = fmt.Sprintf(": \"%v\"", x)
        case nil:
            // Unrecognised object.
            msg = fmt.Sprint(": [error == nil]")
        default:
            // Unrecognised object.
            msg = fmt.Sprintf(": [Unrecognized error] \"%v\"", x)
        }   // End of switch.

        str += fmt.Sprintf("Error %d%s.\n", n, msg)
        n += 1
    }
    return str
}   // End of function Elist::Error.

//=============================================================================
//=============================================================================

/*
Create a new Elist error-message-stack from a given string.
Usage example:
    return elist.New("StructName::MethodName: ErrorDescription");
The return value from elist.New() is of type *Elist, which is assigned to an
"error" interface.
*/
func New(s string) error {
    //------------------//
    //        New       //
    //------------------//
    p := new(Elist)
    // This will never happen.
    if p == nil {
        return nil
    }
    p.value = s
    return p
}   // End of function New.

/*
Create a new Elist error-message-stack from a formatted string.
Usage example:
    return elist.Newf("StructName::MethodName: Error in item %d", n);
The return value from elist.Newf() is of type *Elist, which is assigned to an
"error" interface.
*/
func Newf(format string, args ...interface{}) error {
    //------------------//
    //        Newf      //
    //------------------//
    /*------------------------------------------------------------------------------
        // The slow way to do it!
        n_args := len(args);
        argsCopy := make([]interface{}, n_args);
        for i, arg := range args {
            argsCopy[i] = arg;
            }
        return New(fmt.Sprintf(format, argsCopy...));
    ------------------------------------------------------------------------------*/
    // The quick way to do it!!
    return New(fmt.Sprintf(format, args...))
}   // End of function Newf.

/*
Return a newly created Elist error-message-stack with the new message s at the
head of the stack.

If the argument e is of type *Elist, then the given message s is pushed onto the
given stack e. But if e is any other kind of error, a new Elist object is
created and made to point to e. Thus the new message s is pushed onto either the
old error-message-stack or the old error (of some other type). When
Elist::Error() is called, the new message s is printed before the message or
messages returned by e.

Note if the old error e is of type *Elist, then e itself is used, not a copy of
e.
So the old error-message-stack should not be modified after being linked by this
Push() function.

Usage example:
    var E error = function1();
    return elist.Push(E, "StructName::MethodName: ErrorDescription");
The return value from elist.Push() is of type *Elist, which is assigned to an
"error" interface.
*/
func Push(e error, s string) error {
    //------------------//
    //       Push       //
    //------------------//
    p := new(Elist)
    // This will never happen.
    if p == nil {
        return nil
    }
    p.value = s
    // A nil input-error is not an error. It is a feature!
    if e == nil {
        return p
    }
    q, ok := e.(*Elist)
    if !ok {
        q = new(Elist)
        // Extract the string from the old error and use that.
        //        q.value = e.Error();
        // Make a copy of the entire error.
        // It might contain more than the string.
        q.value = e
    }
    p.next = q
    return p
}   // End of function Push.

/*
Formatted version of Push().
Usage example:
    var E error = function1(n);
    return elist.Pushf(E,
        "StructName::MethodName: Error in call to function1(%d)", n);
The return value from elist.Pushf() is of type *Elist, which is assigned to an
"error" interface.
*/
func Pushf(e error, format string, args ...interface{}) error {
    //------------------//
    //       Pushf      //
    //------------------//
    return Push(e, fmt.Sprintf(format, args...))
}   // End of function Pushf.
