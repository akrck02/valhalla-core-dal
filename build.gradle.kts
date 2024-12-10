import java.io.FileInputStream
import java.util.*

// Variables section
val localProperties: Properties = Properties().apply {
    load(FileInputStream(File(rootProject.rootDir, "local.properties")))
}

group = providers.gradleProperty("organization").getOrElse("")
version = providers.gradleProperty("valhalla.core.dal.version").getOrElse("")

val jdkVersion: Int = providers.gradleProperty("jdk.version").getOrElse("").toInt()
val mavenServerName: String = localProperties.getProperty("maven.server.name")
val mavenServerUrl: String = localProperties.getProperty("maven.server.url")
val mavenServerUser: String = localProperties.getProperty("maven.server.user")
val mavenServerPassword: String = localProperties.getProperty("maven.server.password")
val valhallaSdkVersion: String = providers.gradleProperty("valhalla.core.sdk.version").getOrElse("")

// Plugins section
plugins {
    id("org.jetbrains.kotlin.jvm") version "2.0.21"
    id("maven-publish")
}

// Repositories
repositories {
    mavenCentral()
    mavenLocal()
}

// Testing section
tasks.test {
    useJUnit()
}

// Compilation section
kotlin {
    jvmToolchain(jdkVersion)
}

java {
    withSourcesJar()
    withJavadocJar()
}

// Publishing section
publishing {
    repositories {
        mavenLocal()
        maven {
            name = mavenServerName
            url = uri(mavenServerUrl)
            credentials {
                username = mavenServerUser
                password = mavenServerPassword
            }
            isAllowInsecureProtocol = true
        }
    }
    publications.withType<MavenPublication> {
        from(components["java"])
    }
}

// Dependency section
dependencies {

    testImplementation("org.jetbrains.kotlin:kotlin-test")

    // internal
    implementation("org.akrck02:valhalla.core.sdk:$valhallaSdkVersion")

    // Kotlin coroutine dependency
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.6.4")

    // MongoDB Kotlin driver dependency
    implementation("org.mongodb:mongodb-driver-kotlin-coroutine:4.10.1")

    // Logging
    implementation("org.slf4j:slf4j-api:2.0.0")
    implementation("org.slf4j:slf4j-simple:2.0.0")

    // Env file processing
    implementation("io.github.cdimascio:dotenv-kotlin:6.4.2")

}
