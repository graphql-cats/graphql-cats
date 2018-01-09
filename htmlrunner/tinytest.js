"use strict";
/**
 * Very simple in-browser unit-test library, with zero deps.
 *
 * Background turns green if all tests pass, otherwise red.
 * View the JavaScript console to see failure reasons.
 *
 * Example:
 *
 *   adder.js (code under test)
 *
 *     function add(a, b) {
 *       return a + b;
 *     }
 *
 *   adder-test.html (tests - just open a browser to see results)
 *
 *     <script src="tinytest.js"></script>
 *     <script src="adder.js"></script>
 *     <script>
 *
 *     tests({
 *
 *       'adds numbers': function() {
 *         eq(6, add(2, 4));
 *         eq(6.6, add(2.6, 4));
 *       },
 *
 *       'subtracts numbers': function() {
 *         eq(-2, add(2, -4));
 *       },
 *
 *     });
 *     </script>
 *
 * That's it. Stop using over complicated frameworks that get in your way.
 *
 * -Joe Walnes
 * MIT License. See https://github.com/joewalnes/jstinytest/
 */
var TinyTest = {

    logger: function(logger, style) {
      var pre = document.createElement("pre");
      pre.className = style;

      // script tag exists prior to the body.
      document.addEventListener("DOMContentLoaded", function(e) {
        var body = document.getElementsByTagName("body")[0];
        body.appendChild(pre);
      });

      return function(...msg) {
        logger(msg);
        var txt = document.createTextNode(msg.join("\n") + "\n");
        pre.appendChild(txt);
      };
    },

    run: function(tests) {
        document.write('<style>.fail {background-color:#ff9999;} .pass {background-color:#99ff99;}</style>');
        var error = logger(console.error, 'fail');
        var info = logger(console.log, 'pass');
        var failures = 0;
        for (var testName in tests) {
            var testAction = tests[testName];
            try {
                testAction.apply(this);
                info('Test:', testName, 'OK');
            } catch (e) {
                failures++;
                error('Test:', testName, 'FAILED', e);
                error(e.stack);
            }
        }
    },

    fail: function(msg) {
        throw new Error('fail(): ' + msg);
    },

    assert: function(value, msg) {
        if (!value) {
            throw new Error('assert(): ' + msg);
        }
    },

    assertEquals: function(expected, actual) {
        if (expected != actual) {
            throw new Error('assertEquals() "' + expected + '" != "' + actual + '"');
        }
    },

    assertStrictEquals: function(expected, actual) {
        if (expected !== actual) {
            throw new Error('assertStrictEquals() "' + expected + '" !== "' + actual + '"');
        }
    },

};

var fail               = TinyTest.fail.bind(TinyTest),
    assert             = TinyTest.assert.bind(TinyTest),
    assertEquals       = TinyTest.assertEquals.bind(TinyTest),
    eq                 = TinyTest.assertEquals.bind(TinyTest), // alias for assertEquals
    assertStrictEquals = TinyTest.assertStrictEquals.bind(TinyTest),
    tests              = TinyTest.run.bind(TinyTest),
    logger             = TinyTest.logger.bind(TinyTest);
