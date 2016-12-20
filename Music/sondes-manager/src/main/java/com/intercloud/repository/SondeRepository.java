package com.intercloud.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.intercloud.model.Sonde;

@Repository
public interface SondeRepository extends JpaRepository<Sonde, Long>{

}
