package org.valhalla.core.dal.tool

import org.valhalla.core.sdk.model.exception.ServiceException
import kotlin.test.fail


/** Assert if the given function throws the given [ServiceException]. */
suspend fun assertThrowsServiceException(
    expected: ServiceException,
    code: suspend () -> Unit
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

/** Get if the [ServiceException] status, code and message matches. */
private fun areDifferentExceptions(expected: ServiceException, actual: ServiceException) =
    (expected.status == actual.status)
        .and(expected.message == actual.message)
        .and(expected.status == actual.status)
        .not()