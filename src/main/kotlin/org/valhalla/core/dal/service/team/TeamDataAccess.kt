package org.valhalla.core.dal.service.team

import org.valhalla.core.sdk.model.team.Team
import org.valhalla.core.sdk.repository.TeamRepository

class TeamDataAccess : TeamRepository {
    override suspend fun register(team: Team?): String {
        TODO("Not yet implemented")
    }

    override suspend fun addMember(id: String, userId: String?) {
        TODO("Not yet implemented")
    }

    override suspend fun delete(id: String?) {
        TODO("Not yet implemented")
    }

    override suspend fun deleteMember(id: String, userId: String?) {
        TODO("Not yet implemented")
    }

    override suspend fun get(id: String?): Team {
        TODO("Not yet implemented")
    }

    override suspend fun getAllByMember(userId: String?): List<Team> {
        TODO("Not yet implemented")
    }

    override suspend fun getAllByOwner(userId: String?): List<Team> {
        TODO("Not yet implemented")
    }

    override suspend fun searchByName(filter: String?): List<Team> {
        TODO("Not yet implemented")
    }

    override suspend fun update(id: String?, team: Team?) {
        TODO("Not yet implemented")
    }

    override suspend fun updateOwner(id: String?, userId: String?) {
        TODO("Not yet implemented")
    }
}