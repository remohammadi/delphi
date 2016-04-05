function getUrlVars() {
    var vars = [], hash;
    var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
    for(var i = 0; i < hashes.length; i++)
    {
        hash = hashes[i].split('=');
        vars.push(hash[0]);
        vars[hash[0]] = hash[1];
    }
    return vars;
}

var qId;
var items;
var orderedItemIndexes;

var shuffle = function(array) {
  var currentIndex = array.length, temporaryValue, randomIndex;

  // While there remain elements to shuffle...
  while (0 !== currentIndex) {

    // Pick a remaining element...
    randomIndex = Math.floor(Math.random() * currentIndex);
    currentIndex -= 1;

    // And swap it with the current element.
    temporaryValue = array[currentIndex];
    array[currentIndex] = array[randomIndex];
    array[randomIndex] = temporaryValue;
  }

  return array;
};

var ranks = {};
// rank == 9999 -> user has selected not to vote on v.solution_id
// rank == 9998 -> the solution is new and the user hasn't voted on it yet

var sorting = false;
var comp_tb = {};
var comp_stack = new Array();
var qs_stack = new Array();
var n = 0;

Pair = function(first, second) {
  this.first = first;
  this.second = second;
  this.toString = function() {
    return "<" + first + "/" + second + ">";
  };
  this.reverse = function() {
    return new Pair(this.second, this.first);
  };
};

var add_comp = function(pair) {
  $("#item1").click(function(){
    $("#items").hide();
    var revp = pair.reverse();
    comp_tb[pair.toString()] = -1;
    comp_tb[revp.toString()] = 1;
    $("#items").hide("fast", quicksort_iterate);
    return false;
  });
  // $("#item1").text(items[orderedItemIndexes[pair.first]].title);
  $("#item1").text(orderedItemIndexes[pair.first]);

  $("#item2").click(function(){
    $("#items").hide();
    var revp = pair.reverse();
    comp_tb[pair.toString()] = 1;
    comp_tb[revp.toString()] = -1;
    $("#items").hide("fast", quicksort_iterate);
    return false;
  });
  // $("#item2").text(items[orderedItemIndexes[pair.second]].title);
  $("#item2").text(orderedItemIndexes[pair.second]);

  $("#items").show(300);
}

var swapOII = function(ind1, ind2) {
  var tmp = orderedItemIndexes[ind1];
  orderedItemIndexes[ind1] = orderedItemIndexes[ind2];
  orderedItemIndexes[ind2] = tmp;
}

var quicksort_iterate = function() {
  if (comp_stack.length > 0) {
    add_comp(comp_stack.pop());
    return;
  }

  console.log("-------------------");
  console.log(orderedItemIndexes);
  console.log("comp_tb:");
  console.log(comp_tb);

  var pair = qs_stack.pop();
  console.log("pair:");
  console.log(pair);

  var left = pair.first;
  var right = pair.second;
  var pivotNewIndex = left;
  for (var i = left; i < right; i++) {
    var p = new Pair(right, i);
    if (comp_tb[p.toString()] == 1) {
        swapOII(i, pivotNewIndex);
        pivotNewIndex += 1;
    } else if (!comp_tb[p.toString()]) {
      console.log("Not found: " + p);
    }
  }
  swapOII(pivotNewIndex, right); // Move pivot to its final place
  if ((pivotNewIndex - 1) > left)
      qs_stack.push(new Pair(left, pivotNewIndex - 1));
  if ((pivotNewIndex + 1) < right)
      qs_stack.push(new Pair(pivotNewIndex + 1, right));

  if ((comp_stack.length == 0) && (qs_stack.length > 0)) {
      var pair = qs_stack[qs_stack.length - 1];
      quicksort(pair);
      return;
  } 

  if ((comp_stack.length == 0) && (qs_stack.length == 0)) { // Sort finished
    for (i in orderedItemIndexes) {
        if (ranks[orderedItemIndexes[i]] != 9999)
            ranks[orderedItemIndexes[i]] = i + 1;
    }

    console.log(orderedItemIndexes);
    console.log("Ranks:");
    console.log(ranks);
    console.log("comp_tb:");
    console.log(comp_tb);
    console.log("End");

    // fill_vote_items();
    // $.ajax({url: "{{= URL(r=request, f='save_vote', args=problem.id) }}",
    //     type: "post",
    //     dataType: "json",
    //     data: ranks,
    //     success: update_status,
    //     error: error
    // });
    sorting = false;
  }
};

var quicksort = function(p) {
  var left = p.first;
  var right = p.second;

  console.log("quicksort(" + left + ", " + right + ")");

  comp_tb = {};
  for (var i = right-1; i >= left; i--)
      comp_stack.push(new Pair(i, right));

  console.log("comp_stack");
  console.log(comp_stack);

  add_comp(comp_stack.pop());
}



$(document).ready(function(){
  $("#items").hide();
  var urlVars = getUrlVars();
  qId = urlVars["question"];
  if (! qId) {
    $("#title").text("هنوز امکان ایجاد سؤال جدید وجود ندارد");
    return;
  }
  $.getJSON("/api/questions/" + qId, function(result, status){
    if (status != "success") {
      $("#title").text("سؤال مورد نظر یافت نشد");
      return;
    }
    $("#title").text(result["title"]);
    $("#compareQuestion").text(result["compareQuestion"]);
    items = result["items"];

    orderedItemIndexes = Array.apply(null, items).map(function (_, i) {return i;});
    ranks = Array.apply(null, orderedItemIndexes).map(function (_, i) {return 9998;});
    shuffle(orderedItemIndexes);
    console.log("Start");
    console.log(orderedItemIndexes);

    sorting = true;
    comp_stack = new Array();
    qs_stack = new Array();

    for(n=0; (n < orderedItemIndexes.length) && (ranks[orderedItemIndexes[n]] < 9999); n++);
    var p = new Pair(0, n-1)
    qs_stack.push(p);
    quicksort(p);
  });


});
