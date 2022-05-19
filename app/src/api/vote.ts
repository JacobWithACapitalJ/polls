/** 
 * Cast votes
 * @param pollId {string} - the poll id of the poll the user is voting for
 * @param vote {string} - the vote the user wishes to cast
 */
export default async function castVote(pollId: string, vote: string): Promise<boolean> {
  const endpoint = "https://vote.jomy.dev/api/vote";
  
  const data = JSON.stringify({
    pollId: pollId,
    vote: vote.trim()
  });
  
  console.log("voting for", vote);
  
  let response = await fetch(endpoint, {
    method: 'POST',
    mode: 'cors',
    body: data
  });
  
  console.log(response);
  
  if (!response.ok) {
    throw response.statusText;
  } else {
    return true;
  }
}