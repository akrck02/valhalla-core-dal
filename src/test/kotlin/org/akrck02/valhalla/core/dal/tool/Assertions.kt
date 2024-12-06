package org.akrck02.valhalla.core.dal.tool

import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import kotlin.test.fail

/**
 * Assert if the given code throws the given service exception
 * @param expected The expected exception
 * @param code The code to test
 */
fun assertThrowsServiceException(
    expected: ServiceException,
    code: () -> Unit
) {
    try {
        code()
        fail("Expecting $expected \nBut no exception was raised.\n")
    } catch (actual: ServiceException) {

        if (areDifferentExceptions(expected, actual))
            fail("Expecting $expected \nFound $actual\n")

        println("Raised expected | $actual")

    } catch (unexpected: Exception) {
        fail("Expecting  $expected \nFound $unexpected\n")
    }
}

/**
 * Get if the exception status, code and message matches
 */
private fun areDifferentExceptions(expected: ServiceException, actual: ServiceException) =
    (expected.status == actual.status).and(expected.message == actual.message).and(expected.status == actual.status)
        .not()