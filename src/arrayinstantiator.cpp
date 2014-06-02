/*
 * arrayinstantitator.cpp
 *
 *  Created on: May 21, 2014
 *      Author: clemieux
 */

#include "arrayinstantiator.h"
#include <math.h>

namespace texada {

/**
 * Create instantiator, and fill the array of instantiations
 * @param events events in all traces inputted
 * @param ltlevents set of atomic propositions to be replaced
 */
array_instantiator::array_instantiator(std::set<std::string>* events_,
		spot::ltl::atomic_prop_set *ltlevents) :
			formula_vars(ltlevents), events(events_){
	//setting up some things
	length = formula_vars->size();
	int k = events->size();
	int array_size = pow(k,length);
	return_array = new inst_fxn[array_size];

	// We must produce all permutations with replacement of length length
	// of the events. This can be done by recursion, but here it is done
	// by iteration, going through the entire return array multiple times
	// in order to construct all permutations. Each pass adds a "level"
	// of mapping to each map. For the simple example with two bindings
	// x and y and three events a, b, c, the process goes like this:
	// INITIAL ARRAY:
	// |____|____|____|____|____|____|____|____|____|
	// FIRST PASS, i=0:
	// |x->a|x->b|x->c|x->a|x->b|x->c|x->a|x->b|x->c|
	// SECOND PASS, i=1:
	// |x->a|x->b|x->c|x->a|x->b|x->c|x->a|x->b|x->c|
	// |y->a|y->a|y->a|y->b|y->b|y->b|y->c|y->c|y->c|
	// at this point, the iteration terminates.

	// i is the "level" of the pass, i.e. how deep in formula_vars we are
	int i = 0;
	for (std::set<const spot::ltl::atomic_prop*>::iterator
			formula_it=formula_vars->begin(); formula_it!=formula_vars->end();
			++formula_it){
		//we really should never get to here, as formula_vars ends at length
		if (i >= length) break;
		// since inst_fxn has a string->string map, we obtain the name
		// of the atomic_prop the iterator is pointing to.
		const spot::ltl::atomic_prop* event_var = *formula_it;
		std::string name = event_var->name();
		traverse_and_fill(name, i, k);
		i++;
	}

}

array_instantiator::~array_instantiator() {
	// TODO Auto-generated destructor stub
	// TODO I probably shouldn't delete that new array since it's intended
	// to be used??
}
/**
 * Traverse and fill the inst_fxn array with the correct mapping, according
 * to how "deep" we are in iteration through the formula variables
 * @param formula_event Name of the event to be the first argument of map
 * @param i depth of iteration through formula variables.
 * @param k number of events
 */
void array_instantiator::traverse_and_fill(std::string formula_event, int i,
		int k){
	//indicates when to continue iteration through events
	int switch_var = pow(k,i);
	// find array size
	length = formula_vars->size();
	int array_size = pow(k,length);
	//begin with the first event
	std::set<std::string>::iterator event_iterator = events->begin();
	for (int j = 0; j < array_size ; j++){
		// if we are at switching point, iterate to next event
		if (switch_var == 0){
			++event_iterator;
			// if we've finished passing through the events, go back to the start
			if (event_iterator == events->end()) event_iterator = events->begin();
			// and reset the switch to the top value
			switch_var = pow(k,i);
		}
		return_array[j].inst_map.insert(std::pair<std::string,
				std::string>(formula_event,*event_iterator));
		switch_var--;

	}
}

array_instantiator::inst_fxn* array_instantiator::return_instantiations(){
	return return_array;
}

} /* namespace texada */