var app = (function(){
	var exported = {};
	var codeElem,
	 inputElem,
	 runCodeElem,
	 stackOutput,
	 errElem,
	 pasteElem,
	 pasteNameElem,
	 pasteBtn;

	exported.load = function() {
		codeElem = elmID("current-code");
		inputElem = elmID("repl");
		runCodeElem = elmID("run-code");
		stackOutput = elmID("stack-output");
		errElem = elmID("error-output");
		pasteNameElem = elmID("paste-name");
		pasteBtn = elmID("paste-button");
		pasteElem = elmID("paste-inputs");

		// Clear loading indicators
		codeElem.className = "";
		runCodeElem.className = "";
		pasteBtn.className = "";
		runCodeElem.value = "Run Code";

		pasteBtn.onclick = function(event) {
			exported.addPasteInput(pasteNameElem.value);
			pasteNameElem.value = "";
		}

		runCodeElem.onclick = function() {
			runCode(codeElem.value);
			return false;
		};

		inputElem.onkeypress = function(event) {
			if (event.keyCode != 13 /*Enter*/) {
				return;
			}
			var code = inputElem.value;
			inputElem.value = "Executing code...";
			runCode(code);
			inputElem.value = "";
			return false;
		};
	}

	function runCode(code) {
		errElem.innerHTML = "";
		errElem.className = "hidden";
		try {
			var err = pisc_eval(code);
			if (err) {
				errElem.innerHTML = err.Error();
				errElem.className = "error";

			}
		} catch(ex) {
			errElem.innerHTML = ex.toString();
			errElem.className = "error";
		}
		displayStack();
	}

    function cleanStackDisplay(str) {
        return str.
            replace(/>/g, "&gt;").
            replace(/</g, "&lt;").
            replace(/\n/g, "<br/>");
    }

	function displayStack() {
		var stack = pisc_get_stack(),
			builtHTML = "";
		for (var i = 0; i < stack.length; i++) {
			var dropNum = stack.length - i - 1;
			builtHTML += 
				"<tr><td class='col'>" + cleanStackDisplay(stack[i].String()) + "</td>" + 
					"<td class='col-alt'>&lt;" + stack[i].Type() + "&gt;</td>" +
					'<td class="col"><a href="javascript:app.removeStackElem(' + dropNum + ')">Remove</a></td>' +
				"</tr>";
		}
		stackOutput.innerHTML = builtHTML;
	}

	exported.clearStack = function() {
		pisc_eval("clear-stack");
		displayStack();
	}
	exported.removeStackElem = function(idx) {
		pisc_eval(idx + " pick-del");
		displayStack();
	}

	exported.addPasteInput = function(name) {
		var elem = document.createElement("div");
		elem.innerHTML = 
		'<div>'+
		'<span>'+name+'</span> | ' +
		'<a href="javascript:app.removePasteInput(\''+ name + '\');">Delete</a>' +
		' |' +
		'<span> Use via <code>"'+ name +'" get-paste-text</code> in code/repl </span>' +
		'</div><textarea id="' + name + '-text"cols="50"></textarea>';
		elem.id = name;
		pasteElem.appendChild(elem);
	}

	exported.removePasteInput = function(name) {
		var elem = elmID(name);
		elem.parentNode.removeChild(elem);
	}

	exported.getTextOfPaste = function(name) {
		var elem = elmID(name + "-text");
		return elem.value;
	}

	function elmID(id) {
		return document.getElementById(id);
	}
	return exported;
})();
window.onload = app.load;
