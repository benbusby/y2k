# Y2K Command Cheat Sheet

This provides a quick rundown of how to perform various commands.

All commands follow the same basic flow:

Command ID -> Command Fields -> Command Value

## Command IDs

Every Y2K command should start by referencing a command ID. If you wanted
to print something, your command would start with `9`, if you wanted to
create a new variable, your command would start with `8`, and so on.

<table>
  <tr>
    <th>Command</th>
    <th>Description</th>
    <th>ID</th>
  </tr>
  <tr>
    <td><code>PRINT</code></td>
    <td>Print variable or string</td>
    <td><code>9</code></td>
  </tr>
  <tr>
    <td><code>CREATE</code></td>
    <td>Create a new variable</td>
    <td><code>8</code></td>
  </tr>
  <tr>
    <td><code>MODIFY</code></td>
    <td>Modify an existing variable</td>
    <td><code>7</code></td>
  </tr>
  <tr>
    <td><code>CONDITION</code></td>
    <td>Create a condition</td>
    <td><code>6</code></td>
  </tr>
  <tr>
    <td><code>META</code></td>
    <td>Modify interpreter state</td>
    <td><code>5</code></td>
  </tr>
  <tr>
    <td><code>CONTINUE</code></td>
    <td>Continue (in loop)</td>
    <td><code>4</code></td>
  </tr>
</table>

## Command Fields

After the command ID, the next N digits should complete the required fields
for the specified command. If a command has 3 fields, then the next 3 digits
would be assigned to those fields, for example.

<table>
  <tr>
    <th>Command ID</th>
    <th>Fields</th>
  </tr>
  <tr>
    <td><code>9</code> (<code>PRINT</code>)</td>
    <td>
      <ol>
        <li>Type</li>
        <ul>
          <li>1 --> String</li>
          <li>2 --> Variable</li>
        </ul>
      </ol>
    </td>
  </tr>
  <tr>
    <td><code>8</code> (<code>CREATE</code>)</td>
    <td>
      <ol>
        <li>Variable ID</li>
        <li>Type</li>
        <ul>
          <li>1 --> String</li>
          <li>2 --> Integer</li>
          <li>3 --> Float</li>
          <ul>
            <li>Size should be # digits + 1, with the first digit used for decimal placement.</li>
            <li>Example: <code>3.14</code> would require Size = 4, with the first digit set to <code>1</code> (<code>1314</code>).</li>
          </ul>
          <li>9 --> Copy</li>
        </ul>
        <li>Size</li>
      </ol>
    </td>
  </tr>
  <tr>
    <td><code>7</code> (<code>MODIFY</code>)</td>
    <td>
      <ol>
        <li>Variable ID</li>
        <li>Function</li>
        <ul>
          <li>1 --> <code>+=</code></li>
          <li>2 --> <code>-=</code></li>
          <li>3 --> <code>*=</code></li>
          <li>4 --> <code>/=</code></li>
          <li>5 --> <code>+= other var value</code></li>
          <li>6 --> <code>Copy other var value</code></li>
        </ul>
        <li>Argument Size</li>
      </ol>
    </td>
  </tr>
  <tr>
    <td><code>6</code> (<code>CONDITION</code>)</td>
    <td>
      <ol>
        <li>Variable ID</li>
        <li>Comparison</li>
        <ul>
          <li>1 --> <code>==</code></li>
          <li>2 --> <code><</code></li>
          <li>3 --> <code>></code></li>
          <li>4 --> <code>Is evenly divisible by</code></li>
        </ul>
        <li>Argument Size</li>
      </ol>
    </td>
  </tr>
  <tr>
    <td><code>5</code> (<code>META</code>)</td>
    <td>
      <ol>
        <li>Debug Mode</li>
        <ul>
          <li>0 --> Off</li>
          <li>1 --> On</li>
        </ul>
        <li># of digits</li>
        <ul>
          <li>Updates the number of digits parsed on each pass of the interpreter</li>
        </ul>
      </ol>
    </td>
  </tr>
</table>

## Command Value

After the command ID and fields are set, the next step is to read in a value
that matches the "Size" field (if applicable). For example, if you're creating
a new variable, and the variable size is set to `4`, the next 4 digits would
contain the variable's value. If you're modifying an existing variable, and the
argument size is `3`, the next 3 digits would contain the argument value for the
modifier function.

`PRINT` statements don't have a size field. Depending on the `Type` specified, 
a print statement either:

 - `String`: Converts the following digits into characters until a 2-space
   sequence is reached. See [Character Codes](#character-codes) for help.
 - `Variable`: Use the next digit as a variable ID, and prints that variable

Once the command value is set, the interpreter returns to the start and looks for
the next command ID to repeat this process over again.

### Character Codes

#### Alphabet

<table>
  <tr>
    <th>Number</th>
    <th>Character</th>
    <th>Number</th>
    <th>Character</th>
  </tr>
  <tr>
    <td>0</td>
    <td><code> </code> (whitespace)</td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>1</td>
    <td>a</td>
    <td>27</td>
    <td>A</td>
  </tr>
  <tr>
    <td>2</td>
    <td>b</td>
    <td>28</td>
    <td>B</td>
  </tr>
  <tr>
    <td>3</td>
    <td>c</td>
    <td>29</td>
    <td>C</td>
  </tr>
  <tr>
    <td>4</td>
    <td>d</td>
    <td>30</td>
    <td>D</td>
  </tr>
  <tr>
    <td>5</td>
    <td>e</td>
    <td>31</td>
    <td>E</td>
  </tr>
  <tr>
    <td>6</td>
    <td>f</td>
    <td>32</td>
    <td>F</td>
  </tr>
  <tr>
    <td>7</td>
    <td>g</td>
    <td>33</td>
    <td>G</td>
  </tr>
  <tr>
    <td>8</td>
    <td>h</td>
    <td>34</td>
    <td>H</td>
  </tr>
  <tr>
    <td>9</td>
    <td>i</td>
    <td>35</td>
    <td>I</td>
  </tr>
  <tr>
    <td>10</td>
    <td>j</td>
    <td>36</td>
    <td>J</td>
  </tr>
  <tr>
    <td>11</td>
    <td>k</td>
    <td>36</td>
    <td>K</td>
  </tr>
  <tr>
    <td>12</td>
    <td>l</td>
    <td>37</td>
    <td>L</td>
  </tr>
  <tr>
    <td>13</td>
    <td>m</td>
    <td>38</td>
    <td>M</td>
  </tr>
  <tr>
    <td>14</td>
    <td>n</td>
    <td>39</td>
    <td>N</td>
  </tr>
  <tr>
    <td>15</td>
    <td>o</td>
    <td>40</td>
    <td>O</td>
  </tr>
  <tr>
    <td>16</td>
    <td>p</td>
    <td>41</td>
    <td>P</td>
  </tr>
  <tr>
    <td>17</td>
    <td>q</td>
    <td>42</td>
    <td>Q</td>
  </tr>
  <tr>
    <td>18</td>
    <td>r</td>
    <td>43</td>
    <td>R</td>
  </tr>
  <tr>
    <td>19</td>
    <td>s</td>
    <td>44</td>
    <td>S</td>
  </tr>
  <tr>
    <td>20</td>
    <td>t</td>
    <td>45</td>
    <td>T</td>
  </tr>
  <tr>
    <td>21</td>
    <td>u</td>
    <td>46</td>
    <td>U</td>
  </tr>
  <tr>
    <td>22</td>
    <td>v</td>
    <td>47</td>
    <td>V</td>
  </tr>
  <tr>
    <td>23</td>
    <td>w</td>
    <td>48</td>
    <td>W</td>
  </tr>
  <tr>
    <td>24</td>
    <td>x</td>
    <td>49</td>
    <td>X</td>
  </tr>
  <tr>
    <td>25</td>
    <td>y</td>
    <td>50</td>
    <td>Y</td>
  </tr>
  <tr>
    <td>26</td>
    <td>z</td>
    <td>51</td>
    <td>Z</td>
  </tr>
</table>

#### Numeric

<table>
  <tr>
    <th>Number</th>
    <th>Character</th>
  </tr>
  <tr>
    <td>52</td>
    <td>1</td>
  </tr>
  <tr>
    <td>53</td>
    <td>2</td>
  </tr>
  <tr>
    <td>54</td>
    <td>3</td>
  </tr>
  <tr>
    <td>55</td>
    <td>4</td>
  </tr>
  <tr>
    <td>56</td>
    <td>5</td>
  </tr>
  <tr>
    <td>57</td>
    <td>6</td>
  </tr>
  <tr>
    <td>58</td>
    <td>7</td>
  </tr>
  <tr>
    <td>59</td>
    <td>8</td>
  </tr>
  <tr>
    <td>60</td>
    <td>9</td>
  </tr>
  <tr>
    <td>61</td>
    <td>0</td>
  </tr>
</table>

#### Symbols

!@#$%^&*()+-<>.,

<table>
  <tr>
    <th>Number</th>
    <th>Character</th>
  </tr>
  <tr>
    <td>62</td>
    <td>!</td>
  </tr>
  <tr>
    <td>63</td>
    <td>@</td>
  </tr>
  <tr>
    <td>64</td>
    <td>#</td>
  </tr>
  <tr>
    <td>65</td>
    <td>$</td>
  </tr>
  <tr>
    <td>66</td>
    <td>%</td>
  </tr>
  <tr>
    <td>67</td>
    <td>^</td>
  </tr>
  <tr>
    <td>68</td>
    <td>&</td>
  </tr>
  <tr>
    <td>69</td>
    <td>*</td>
  </tr>
  <tr>
    <td>70</td>
    <td>(</td>
  </tr>
  <tr>
    <td>71</td>
    <td>)</td>
  </tr>
  <tr>
    <td>72</td>
    <td>+</td>
  </tr>
  <tr>
    <td>73</td>
    <td>-</td>
  </tr>
  <tr>
    <td>74</td>
    <td><</td>
  </tr>
  <tr>
    <td>75</td>
    <td>></td>
  </tr>
  <tr>
    <td>76</td>
    <td>.</td>
  </tr>
  <tr>
    <td>77</td>
    <td>,</td>
  </tr>
</table>
