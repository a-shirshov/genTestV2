export const parseMessage = (text) => {
    const thinkSection = text.match(/<think>(.*?)<\/think>/s);
    return {
      thinkSection: thinkSection?.[1] || '',
      cleanedText: text.replace(/<\/?think>/g, '')
    };
  };

export default parseMessage;