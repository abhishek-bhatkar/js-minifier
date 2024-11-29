/*!
 * Test Library v1.0.0
 * (c) 2024 Test Author
 * Released under the MIT License
 */

/*! Important runtime warning - Do not remove */
const IMPORTANT_CONSTANT = 42;

// Regular comment that should be removed
function someFunction() {
    /* This comment should be removed */
    return IMPORTANT_CONSTANT;
}

/*!
 * This is a multi-line license comment
 * that should be preserved
 * @license
 */
class ImportantClass {
    // This regular comment should be removed
    constructor() {
        /*! Runtime critical comment */
        this.value = IMPORTANT_CONSTANT;
    }
}

/* This is a normal multi-line comment
   that should be removed during
   minification */
const result = someFunction();
