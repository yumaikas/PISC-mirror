-- Errors
local EOF = "End of file"

local skipchars = {
	n = '\n',
	t = '\t',
	r = '\r',
	["\\"] = '\\',
	["[\""] = '"',
}

function string.starts(String,Start)
   return string.sub(String, 1, string.len(Start)) == Start
end

function string.ends(String,End)
   return End == '' or string.sub(String, -string.len(End))== End
end

local CodePosition = {
	Name = "CodePosition",
	Spec = {
		LineNumber = "number",
		Offset = "number",
		Source = "string"
	},
}

local CodeList = {
	Name = "Codelist",
	Spec = {
		CodePosition = "table",
		Idx = "number",
		Code = "string",
		FileName = "string"
	},
}

local CodeQuotation = {
	Name = "CodeQuotation",
	Spec = {
		Idx = "number",
		Words = "table",
		CodePositions = "table"
	},
}

function  _ctor(_type, overrides)
	function default_ctor(__type)
		local value = {}
		for k,v in pairs(__type.Spec) do
			if v == "number" then
				value[k] = 0.0
			elseif v == "string" then
				value[k] = ""
			elseif v == "boolean" then
				value[k] = false
			elseif v == "table" then
				value[k] = {}
			else 
				assert(false, "Type "..v.." isn't a supported defaultable type")
			end
		end
		return value
	end
	local value = default_ctor(_type)
	if overrides then
		for k,v in pairs(overrides) do
			value[k] = v
		end
	end
	return value
end

function _type(value, __type)
	for k,v in pairs(__type.Spec) do
		if type(value[k]) ~= v then
			return false, __type.Name .. "'s " .. k .. " isn't a " .. v .. " but a " .. type(value[k])
		end
	end
	return true, ""
end

function make_lexer(code, position) 
	assert(_type(position, CodePosition))
	local basis = _ctor(CodeList, {Code = code, CodePosition = position} )
	assert(_type(basis, CodeList))

	basis.nextWord = function (self)
		local currentWord = ""
		local skipChar = false
		local inLineComment = false
		currentLine = ""

		if self.Idx > #self.Code then
			return {str = currentWord}, EOF
		end

		for v in self.Code:sub(self.Idx):gmatch(".") do
			-- Emit these words as they are found
			if currentWord == "${"
				or currentWord == "{"
				or currentWord == "}" 
				or currentWord == "["
			    or currentWord == "]" then
			    return {str = currentWord}, nil
			end

			self.Idx = self.Idx + #v
			self.CodePosition.Offset = self.CodePosition.Offset + #v

			if v == "\n" then
				currentLine = ""
				self.CodePosition.LineNumber = self.CodePosition.LineNumber + 1
				self.CodePosition.Offset = 0
			end

			-- This is the logic for handling [] and {} being able to be adjacent to other words.
			if not(inString or inLineComment) then
				if v == "{" and currentWord:ends("$") and currentWord ~= "$" then
					-- Unread both the { and $
					c.Idx = c.Idx - #v - #"$"
					c.Offset = c.Offset - #v - #"$"
					return {str = currentWord:sub(1, -(#"$" + 1))}, nil
				end
				if v == "{" and currentWord:ends("$") and currentWord ~= "$" then
					return {str =  "${"}, nil
				end
				if (v == "[" or v == "{" or v == "}" or v == "]") and #currentWord > 0 then
					-- Unread the {,},[,]
					self.Idx = self.Idx - #v
					self.CodePosition.Offset = self.CodePosition.Offset - #v
					return {str = currentWord}, nil
				end
			end
			currentLine = currentLine .. v
			if inLineComment then
				if v == '\n' or v == '\r' then
					-- TODO: Mark these as comments somehow?
					return {str = currentWord}, nil
				else
					currentWord = currentWord .. v
				end
			end

			-- Handle \n, \t, \r \\ and \"
			if skipChar and skipchars[v] then
				if skipchars[v] ~= nil then
					currentWord = currentWord .. skipchars[v]
				else
					return nil, 
						"Invalid escape sequence: "..v
						.." current word: "..currentWord
						.." line: "..self.CodePosition.LineNumber
				end
			end

			-- This currently needs to be the terminal logic for this bit of the lexer
			if v == "\\" and inString then
					skipChar = true
			elseif v == "#" then
				if not inString then
					inLineComment = true
				end
				currentWord = currentWord .. v
			elseif v == '"' then
				if inString then
					currentWord =  currentWord .. '"'
					inString = false
				else
					inString = true
					currentWord = currentWord .. v
				end
			elseif v == ' ' or v == '\t' or v == '\n' or v == '\r' then
				if inString then
					currentWord = currentWord .. v
				elseif #currentWord > 0 then
					return {str = currentWord}, nil
				end
			else
				currentWord = currentWord .. v
			end
		end
		if inString then
			return nil, "Unterminated string!"
		end
		return {str = currentWord}, nil
	end
	return basis
end

function stringToQuotation(code, position) -- Quotation, error
	local basis = make_lexer(code, position)
	local quot = _ctor(CodeQuotation)

	local error
	local word

 	while error == nil do
 		word, err = basis:nextWord()
 		if err == EOF then
 			return quot, nil
 		end

 		if err ~= nil then
 			return nil, err
 		end
 		table.insert(quot.Words, word)
 		table.insert(quot.CodePositions, _ctor(CodePosition, basis.CodePosition))
 	end
 	return quot, nil
end

-- Print anything - including nested tables
function table_print (tt, indent, done)
  done = done or {}
  indent = indent or 0
  if type(tt) == "table" then
    for key, value in pairs (tt) do
      io.write(string.rep (" ", indent)) -- indent it
      if type (value) == "table" and not done [value] then
        done [value] = true
        io.write(string.format("[%s] => table\n", tostring (key)));
        io.write(string.rep (" ", indent+4)) -- indent it
        io.write("(\n");
        table_print (value, indent + 7, done)
        io.write(string.rep (" ", indent+4)) -- indent it
        io.write(")\n");
      else
        io.write(string.format("[%s] => %s\n",
            tostring (key), tostring(value)))
      end
    end
  else
    io.write(tt .. "\n")
  end
end

function lex_string(code)
	local pos = _ctor(CodePosition, {Source = "lex_string"})
	return stringToQuotation(code, _ctor(CodePosition))
end

table_print(lex_string('1 2 3 \n\n[ "foo" "bar" ] '))