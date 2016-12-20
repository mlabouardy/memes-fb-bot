package com.intercloud.api;

import java.util.Collection;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

import com.intercloud.model.Sonde;
import com.intercloud.repository.SondeRepository;

@RestController
public class SondeResource {
	
	@Autowired
	private SondeRepository sondeRepository;
	
	@RequestMapping(value="/sondes", method=RequestMethod.GET, produces="application/json")
	public Collection<Sonde> getSondes(){
		return sondeRepository.findAll();
	}
	
	
	@RequestMapping(value="/sondes", method=RequestMethod.POST, consumes="application/json",produces="application/json")
	public void save(@RequestBody Sonde sonde){
		sondeRepository.save(sonde);
	}
}
