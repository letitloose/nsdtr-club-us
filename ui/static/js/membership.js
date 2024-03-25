
membershipTypeSelect = document.getElementById("membership-type");
membershipTypeSelect.addEventListener('change',
function() { updateAmount(this.value); },
false);

membershipAmount = document.getElementById("membership-amount");

function updateAmount(value){
    if(value == "SI"){
        membershipAmount.value = 25;
    }else if(value == "JT"){
        membershipAmount.value = 35;
    }else{
        membershipAmount.value = 12;
    }
    
}

updateAmount("SI");

roster = document.getElementById("roster");
roster.addEventListener('change',
function() { updateRoster(this.checked); },
false);

rosterAmount = document.getElementById("roster-amount");

function updateRoster(checked){
    console.log("derr " + checked)
    if(checked){
        rosterAmount.value = 5;
    }else{
        rosterAmount.value = 0;
    }
}