# Open Gallery Info

  Open Gallery Info is an open data project that aims to facilitate open data
  publishing for art galleries.

## Sepcification

  `Gallery` contains a gallery information with user-defined [UUID].
  `Exhibition` is a [CSV] file that contains a list of exhibition information.

### Gallery

  JSON formatted gallery data.

#### id

  Unique identifier. It MUST be an [UUID] formatted string. It SHOULD be
  generated by a computer software rather than selecting numbers by preference.

#### name

  Gallery name.

#### address

  Gallery location.

#### exhibitions

  Array of string.

#### about optional

  Describe about a gallery here.

#### open_at optional

  Opening hour. e.g. "10:00"

#### close_at optional

  Closeing hour. e.g. "18:00"

#### close_on optional

  holiday.


### Exhibition, CSV

CSV formatted exhibition data.

#### id

  User-defined identifier. It MUST be an unique identifier. It MUST NOT conflict
  with other data entry.

#### title

  Exhibition title.

#### start

  Start date.

#### end

  End date.

#### description, optional

  Description.

#### alert, optional

  Alert of an infomation for an exhibition.

#### note, optional

  Note of an information for an exhibition.

[UUID]: http://en.wikipedia.org/wiki/Universally_unique_identifier
[JSON]: http://en.wikipedia.org/wiki/JSON
[CSV]: http://en.wikipedia.org/wiki/Comma-separated_values
