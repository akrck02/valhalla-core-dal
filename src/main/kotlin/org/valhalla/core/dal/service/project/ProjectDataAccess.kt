package org.valhalla.core.dal.service.project

import org.valhalla.core.sdk.model.project.Project
import org.valhalla.core.sdk.repository.ProjectRepository

class ProjectDataAccess : ProjectRepository {

    override suspend fun delete(id: String?) {
        TODO("Not yet implemented")
    }

    override suspend fun get(id: String?): Project {
        TODO("Not yet implemented")
    }

    override suspend fun getAllByMember(userId: String?): List<Project> {
        TODO("Not yet implemented")
    }

    override suspend fun register(project: Project?): String {
        TODO("Not yet implemented")
    }

    override suspend fun update(id: String?, project: Project?) {
        TODO("Not yet implemented")
    }

}

